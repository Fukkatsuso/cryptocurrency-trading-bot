package persistence

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
	"google.golang.org/api/option"
)

// モックデータ
// というよりは，本番DBからエクスポートされた価格データファイルを取得する
// candleTableNameも固定するのでproductCodeとdurationも固定することにする
type candleMockRepository struct {
	candleTableName string
	timeFormat      string
	productCode     string
	duration        time.Duration
	candles         []model.Candle
}

func NewCandleMockRepository(candleTableName, timeFormat, productCode string, duration time.Duration) repository.CandleRepository {
	if candleTableName == "" {
		return nil
	}

	cr := &candleMockRepository{
		candleTableName: candleTableName,
		timeFormat:      timeFormat,
		productCode:     productCode,
		duration:        duration,
	}

	candles, err := cr.fetchAll(productCode, duration)
	if err != nil {
		return nil
	}
	cr.candles = candles

	return cr
}

const (
	dataDir = "/mock_data"
	saKey   = "/sa_key"
)

var (
	MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE")
	GCS_BUCKET     = os.Getenv("GCS_BUCKET")
)

func (cr *candleMockRepository) Save(candle model.Candle) error {
	for i, candle := range cr.candles {
		if candle.Time().Equal(candle.Time()) {
			cr.candles[i] = candle
			return nil
		}
	}

	cr.candles = append(cr.candles, candle)
	sort.Slice(cr.candles, func(i, j int) bool {
		timeBefore := cr.candles[i].Time().Time()
		timeAfter := cr.candles[j].Time().Time()
		return timeBefore.Before(timeAfter)
	})

	return nil
}

func (cr *candleMockRepository) FindByCandleTime(productCode string, duration time.Duration, timeTime model.CandleTime) (*model.Candle, error) {
	for _, candle := range cr.candles {
		if candle.Time().Equal(timeTime) {
			return &candle, nil
		}
	}

	return nil, nil
}

func (cr *candleMockRepository) FindAll(productCode string, duration time.Duration, limit int64) ([]model.Candle, error) {
	if limit < 0 {
		return cr.candles, nil
	}

	if lenCandles := int64(len(cr.candles)); lenCandles > limit {
		return cr.candles[lenCandles-limit:], nil
	}
	return cr.candles, nil
}

func (cr *candleMockRepository) fetchAll(productCode string, duration time.Duration) ([]model.Candle, error) {
	dirPath := path.Join(dataDir, GCS_BUCKET)
	if !exists(dirPath) {
		os.MkdirAll(dirPath, 0777)
	}
	objectName := fmt.Sprintf("%s.%s.csv", MYSQL_DATABASE, cr.candleTableName)
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
	candles := make([]model.Candle, 0)
	for {
		// time, open, close, high, low, volume
		line, err := reader.Read()
		if err != nil {
			break
		}

		time, err := time.Parse(cr.timeFormat, line[0])
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

		candleTime := model.NewCandleTime(time)
		candle := model.NewCandle(productCode, duration, candleTime, open, close, high, low, volume)
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
