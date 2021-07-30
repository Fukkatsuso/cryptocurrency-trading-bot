package model

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"google.golang.org/api/option"
)

const (
	dataDir = "/mock_data"
	saKey   = "/sa_key"
)

var (
	MYSQL_DATABASE  = os.Getenv("MYSQL_DATABASE")
	CANDLE_TABLE    = config.CandleTableName
	PRODUCT_CODE    = config.ProductCode
	CANDLE_DURATION = config.CandleDuration
	GCS_BUCKET      = os.Getenv("GCS_BUCKET")
)

func CandleMockData() ([]Candle, error) {
	dirPath := path.Join(dataDir, GCS_BUCKET)
	if !exists(dirPath) {
		os.MkdirAll(dirPath, 0777)
	}
	objectName := fmt.Sprintf("%s.%s.csv", MYSQL_DATABASE, CANDLE_TABLE)
	filePath := path.Join(dirPath, objectName)
	if !exists(filePath) {
		downloadGCSObject(GCS_BUCKET, objectName, filePath)
	}

	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	candles := make([]Candle, 0)
	for {
		// time, open, close, high, low, volume
		line, err := reader.Read()
		if err != nil {
			break
		}

		time, err := time.Parse(config.TimeFormat, line[0])
		if err != nil {
			return nil, err
		}

		open, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}

		close, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, err
		}

		high, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return nil, err
		}

		low, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			return nil, err
		}

		volume, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			return nil, err
		}

		candle := NewCandle(PRODUCT_CODE, CANDLE_DURATION, time, open, close, high, low, volume)
		candles = append(candles, *candle)
	}

	return candles, nil
}

func exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func downloadGCSObject(bucketName, objectName, localFilePath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(saKey))
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

	err = ioutil.WriteFile(localFilePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
