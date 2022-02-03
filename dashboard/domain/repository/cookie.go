package repository

import "net/http"

type Cookie interface {
	Set(w http.ResponseWriter, value map[string]string) error
	GetValue(r *http.Request) (map[string]string, error)
	Delete(w http.ResponseWriter)
}
