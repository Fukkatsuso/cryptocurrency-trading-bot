package cloudfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

var (
	GCP_PROJECT       = os.Getenv("GCP_PROJECT")
	CLOUDSQL_INSTANCE = os.Getenv("CLOUDSQL_INSTANCE")
	DATABASE          = os.Getenv("MYSQL_DATABASE")
	GCS_BUCKET        = os.Getenv("GCS_BUCKET")
)

const (
	CandleTableName = "eth_candles"
)

func ExportDatabaseToStorage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	c, err := google.DefaultClient(ctx, sqladmin.CloudPlatformScope)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqladminService, err := sqladmin.New(c)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// "gs://bucketName/fileName"
	storageObject := fmt.Sprintf("%s.%s.csv", DATABASE, CandleTableName)
	uri := fmt.Sprintf("gs://%s/%s", GCS_BUCKET, storageObject)
	columns := "time, open, close, high, low, volume"
	selectQuery := fmt.Sprintf("SELECT %s FROM %s.%s ORDER BY time ASC", columns, DATABASE, CandleTableName)
	rb := &sqladmin.InstancesExportRequest{
		ExportContext: &sqladmin.ExportContext{
			Kind:     "sql#exportContext",
			FileType: "CSV",
			Offload:  true,
			Uri:      uri,
			CsvExportOptions: &sqladmin.ExportContextCsvExportOptions{
				SelectQuery: selectQuery,
			},
		},
	}
	fmt.Println("uri:", uri)
	fmt.Println("selectQuery:", selectQuery)

	// オブジェクトは上書きできないため，削除後に新規作成（エクスポート）する
	// なるべくエクスポートの直前に削除して，オブジェクトがない空白時間を減らす
	err = deleteStorageObject(ctx, GCS_BUCKET, storageObject)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// エクスポート実行
	resp, err := sqladminService.Instances.Export(GCP_PROJECT, CLOUDSQL_INSTANCE, rb).Context(ctx).Do()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Error != nil {
		fmt.Printf("%#v\n", resp.Error)
		js, err := json.Marshal(resp.Error)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, string(js), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, uri)
}

func deleteStorageObject(ctx context.Context, bucketName, objectName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// 削除してみる
	err = client.Bucket(bucketName).Object(objectName).Delete(ctx)
	// "存在しないオブジェクトである"というエラーなら問題ない
	// 逆にこれ以外のエラーはエラーとして扱う
	if err != storage.ErrObjectNotExist {
		return err
	}

	return nil
}
