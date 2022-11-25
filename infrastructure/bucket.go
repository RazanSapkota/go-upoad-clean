package infrastructure

import (
	"context"
	"example/go-api/lib"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// NewBucketStorage creates a new storage client
func NewBucketStorage(env *lib.Env) *storage.Client {
	bucketName := env.StorageBucketName
	ctx := context.Background()
	if bucketName == "" {
		log.Fatal("Please check your env file for StorageBucketName")
	}
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		log.Fatalf("Provided bucket %v doesn't exists", bucketName)
	}
	if err != nil {
		log.Fatalf("Cloud bucket error: %v", err.Error())
	}
	return client
}
