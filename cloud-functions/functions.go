package cloudfunctions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	uri := fmt.Sprintf("gs://%s/%s.%s.csv", GCS_BUCKET, DATABASE, CandleTableName)
	databases := []string{
		DATABASE,
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY time ASC", CandleTableName)
	rb := &sqladmin.InstancesExportRequest{
		ExportContext: &sqladmin.ExportContext{
			Kind:      "sql#exportContext",
			FileType:  "CSV",
			Uri:       uri,
			Databases: databases,
			CsvExportOptions: &sqladmin.ExportContextCsvExportOptions{
				SelectQuery: selectQuery,
			},
		},
	}

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
}
