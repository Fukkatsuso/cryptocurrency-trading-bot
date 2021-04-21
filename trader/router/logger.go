package router

import (
	"fmt"
	"net/http"
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		fmt.Printf("[%s] %s\n", method, path)
		h.ServeHTTP(w, r)
	})
}
