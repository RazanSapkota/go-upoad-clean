package controllers

import (
	"bytes"
	"example/go-api/service"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type UploadController interface {
	Upload(ctx *gin.Context)
}

type uploadController struct {
	jwtService   service.JWTService
	bucket       service.BucketService
}

func NewUploadController(jwtService service.JWTService, bucketService service.BucketService) UploadController {
	return &uploadController{
		jwtService: jwtService,
		bucket:     bucketService,
	}
}

func (controller *uploadController) 	Upload(c *gin.Context){
	_, ctx := errgroup.WithContext(c.Request.Context())
	file, fileHeader, err := c.Request.FormFile("file")
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":err.Error(),
		}) 
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fileByte, err := io.ReadAll(file)
	
		if err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"message":err.Error(),
			}) 
			return
		}
	fileReader := bytes.NewReader(fileByte)
	randUUID, _ := uuid.NewRandom()
	fileName := randUUID.String() + ext
	urlResponse, err:=controller.bucket.UploadFile(ctx,fileReader,fileName,fileHeader.Filename)
	if err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message":err.Error(),
		}) 
		return
	}

	c.JSON(200, gin.H{
	   "message":urlResponse,
   }) 
   
   }