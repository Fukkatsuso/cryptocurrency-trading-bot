package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().In(config.LocalTime).Format("2006-01-02 15:04:05")
		method := r.Method
		path := r.URL.Path
		fmt.Printf("%s [%s] %s", now, method, path)
		h.ServeHTTP(w, r)
	})
}
