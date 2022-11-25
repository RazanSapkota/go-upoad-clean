package service

import (
	"context"
	"example/go-api/lib"
	"io"
	"log"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
)

type BucketService struct {
	client *storage.Client
	env    *lib.Env
}

func NewBucketService(
	client *storage.Client,
	env *lib.Env,
) BucketService {
	return BucketService{
		client: client,
		env:    env,
	}
}

func (s *BucketService) UploadFile(ctx context.Context, file io.Reader, fileName string, originalFileName string) (string, error) {
	var bucketName = s.env.StorageBucketName

	if bucketName == "" {
		log.Fatal("No bucket name in env.")
	}
	_, err := s.client.Bucket(bucketName).Attrs(ctx)

	if err == storage.ErrBucketNotExist {
		log.Fatalf("provided bucket %v doesn't exists", bucketName)
	}
	if err != nil {
		log.Fatalf("cloud bucket error: %v", err.Error())
	}

	wc := s.client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	//wc.ContentType = "application/octet-stream"
	//wc.ContentDisposition = "attachment; filename=" + originalFileName

	_, err = io.Copy(wc, file)
	if err != nil {
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}

	u, err := url.ParseRequestURI("/" + bucketName + "/" + wc.Attrs().Name)

	if err != nil {
		return "", err
	}

	path := u.EscapedPath()
	path = strings.Replace(path, "/"+bucketName, "", 1)
	path = strings.Replace(path, "/", "", 1)

	return path, nil
}
