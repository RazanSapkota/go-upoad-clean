package controllers

import (
	"bytes"
	"context"
	"errors"
	"example/go-api/lib"
	"example/go-api/service"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"golang.org/x/sync/errgroup"
)

type UploadController interface {
	Upload(ctx *gin.Context)
	uploadFile(
		ctx context.Context,
		errGroup *errgroup.Group,
		conf UploadConfig,
		file multipart.File,
		fileHeader *multipart.FileHeader,
		uploadedFiles *lib.UploadedFiles,
	) error
}

type Extension string

const (
	JPEGFile Extension = ".jpeg"
	JPGFile  Extension = ".jpg"
	PNGFile  Extension = ".png"
)

type UploadConfig struct {
	// FieldName where to pull multipart files from
	FieldName string

	// BucketFolder where to put the uploaded files to
	BucketFolder string

	// Extensions array of extensions
	Extensions []Extension

	// ThumbnailEnabled set whether thumbnail is enabled or nor
	ThumbnailEnabled bool

	// ThumbnailWidth set thumbnail width
	ThumbnailWidth uint

	// WebpEnabled set whether thumbnail is enabled or nor
	WebpEnabled bool

	// Multiple set whether to upload multiple files with same key name
	Multiple bool
}
type uploadController struct {
	jwtService service.JWTService
	bucket     service.BucketService
	config     UploadConfig
}

func (u *uploadController) Config() UploadConfig {
	return UploadConfig{
		FieldName:        "file",
		BucketFolder:     "",
		Extensions:       []Extension{JPEGFile, PNGFile, JPGFile},
		ThumbnailEnabled: true,
		ThumbnailWidth:   100,
		Multiple:         false,
	}
}

func NewUploadController(jwtService service.JWTService, bucketService service.BucketService) UploadController {

	return &uploadController{
		jwtService: jwtService,
		bucket:     bucketService,
		config: UploadConfig{
			FieldName:        "file",
			BucketFolder:     "",
			Extensions:       []Extension{JPEGFile, PNGFile, JPGFile},
			ThumbnailEnabled: false,
			ThumbnailWidth:   100,
			Multiple:         false,
			WebpEnabled:      true,
		},
	}
}

func (controller *uploadController) Upload(c *gin.Context) {
	errGroup, ctx := errgroup.WithContext(c.Request.Context())

	var uploadedFiles lib.UploadedFiles
	conf := controller.config
	if conf.Multiple {
		form, _ := c.MultipartForm()
		files := form.File[conf.FieldName]

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			defer file.Close()
			err = controller.uploadFile(ctx, errGroup, conf, file, fileHeader, &uploadedFiles)
			if err != nil {
				log.Println("file-upload-error: ", err.Error())
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				c.Abort()
				return
			}
		}
	} else {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		err = controller.uploadFile(ctx, errGroup, conf, file, fileHeader, &uploadedFiles)
		if err != nil {
			log.Println("file-upload-error: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
	}

	if err := errGroup.Wait(); err != nil {
		log.Println("file-upload-error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"uploadedFiles": uploadedFiles,
	})

}
func (controller uploadController) uploadFile(
	ctx context.Context,
	errGroup *errgroup.Group,
	conf UploadConfig,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	uploadedFiles *lib.UploadedFiles,
) error {
	if file == nil || fileHeader == nil {
		log.Println("file and fileheader nil value is passed")
		return nil
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !controller.matchesExtension(conf, ext) {
		return errors.New("file extension not supported")
	}
	fileByte, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	uploadFileName, fileUID := controller.randomFileName(conf, ext)

	//OriginalImage
	errGroup.Go(func() error {
		fileReader := bytes.NewReader(fileByte)
		urlResponse, err := controller.bucket.UploadFile(ctx, fileReader, uploadFileName, fileHeader.Filename)
		*uploadedFiles = append(*uploadedFiles, lib.UploadMetadata{
			FieldName: conf.FieldName,
			FileName:  fileHeader.Filename,
			URL:       urlResponse,
			FileUID:   fileUID,
			Size:      fileHeader.Size,
		})
		return err
	})

	//WebPImage
	if conf.WebpEnabled {
		OriginalWebpReader := bytes.NewReader(fileByte)
		errGroup.Go(func() error {
			var webpBuf bytes.Buffer
			img, err := controller.getImage(OriginalWebpReader, ext)
			if err != nil {
				return err
			}
			if err := webp.Encode(&webpBuf, img, &webp.Options{Lossless: true}); err != nil {
				return err
			}
			webpReader := bytes.NewReader(webpBuf.Bytes())
			resizeFileName := controller.bucketPath(conf, fmt.Sprintf("%s_webp%s", fileUID, ".webp"))

			if _, err := controller.bucket.UploadFile(ctx, webpReader, resizeFileName, strings.ReplaceAll(fileHeader.Filename, ext, "")+".webp"); err != nil {
				return err
			}
			return nil
		})
	}

	if conf.ThumbnailEnabled {
		thumbReader := bytes.NewReader(fileByte)

		errGroup.Go(func() error {
			var thumbBuf bytes.Buffer
			img, err := controller.getImage(thumbReader, ext)
			if err != nil {
				return err
			}
			resizeImage := resize.Resize(conf.ThumbnailWidth, 0, img, resize.Lanczos3)
			if Extension(ext) == JPGFile || Extension(ext) == JPEGFile {
				if err := jpeg.Encode(&thumbBuf, resizeImage, nil); err != nil {
					return err
				}
			}
			if Extension(ext) == PNGFile {
				if err := png.Encode(&thumbBuf, resizeImage); err != nil {
					return err
				}
			}
			newThumbReader := bytes.NewReader(thumbBuf.Bytes())
			resizeFileName := controller.bucketPath(conf, fmt.Sprintf("%s_thumb%s", fileUID, ext))
			_, err = controller.bucket.UploadFile(ctx, newThumbReader, resizeFileName, fileHeader.Filename)
			if err != nil {
				return err
			}
			return nil
		})

		if conf.WebpEnabled {
			webpReader := bytes.NewReader(fileByte)
			errGroup.Go(func() error {
				var webpBuf bytes.Buffer
				img, err := controller.getImage(webpReader, ext)
				if err != nil {
					return err
				}

				resizeImage := resize.Resize(conf.ThumbnailWidth, 0, img, resize.Lanczos3)
				err = webp.Encode(&webpBuf, resizeImage, &webp.Options{Lossless: true})
				if err != nil {
					return err
				}

				webpReader := bytes.NewReader(webpBuf.Bytes())
				resizeFileName := controller.bucketPath(conf, fmt.Sprintf("%s_thumb%s", fileUID, ".webp"))

				_, err = controller.bucket.UploadFile(ctx, webpReader, resizeFileName, fileHeader.Filename)
				if err != nil {
					return err
				}

				return nil
			})
		}

	}
	return nil
}

func (u *uploadController) getImage(file io.Reader, ext string) (image.Image, error) {
	if Extension(ext) == JPGFile || Extension(ext) == JPEGFile {
		return jpeg.Decode(file)
	}
	if Extension(ext) == PNGFile {
		return png.Decode(file)
	}
	return nil, errors.New("file extension not supported")
}

func (u *uploadController) randomFileName(c UploadConfig, ext string) (randomName, uid string) {
	randUUID, _ := uuid.NewRandom()
	fileName := randUUID.String() + ext
	return u.bucketPath(c, fileName), randUUID.String()
}

func (u *uploadController) bucketPath(c UploadConfig, fileName string) string {
	if c.BucketFolder != "" {
		return fmt.Sprintf("%s/%s", c.BucketFolder, fileName)
	}
	return fileName
}

func (u *uploadController) matchesExtension(c UploadConfig, ext string) bool {
	for _, e := range c.Extensions {
		if e == Extension(ext) {
			return true
		}
	}
	return false
}
