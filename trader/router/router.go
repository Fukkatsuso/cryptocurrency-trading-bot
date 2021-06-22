package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/controller"
)

func Run() {
	http.HandleFunc("/", controller.HelloWorldHandler)
	http.HandleFunc("/fetch-ticker", controller.FetchTickerHandler)
	http.HandleFunc("/trade", controller.TradeHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	fmt.Printf("listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, logger(http.DefaultServeMux)); err != nil {
		fmt.Println(err)
	}
}
