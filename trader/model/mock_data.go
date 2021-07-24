package model

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"time"

	"cloud.google.com/go/storage"
)

const (
	dataDir = "mock_data"
)

var (
	GCS_BUCKET = os.Getenv("GCS_BUCKET")
)

func downloadGCSObject(bucketName, objectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}

	filePath := path.Join(dataDir, bucketName, objectName)
	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
