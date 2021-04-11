package router

import (
	"log"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/controller"
)

func Run() {
	http.HandleFunc("/", controller.HelloWorldHandler)
	http.HandleFunc("/fetch-board", controller.FetchBoardHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, logger(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
