package main

import (
	"log"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/router"
)

func main() {
	log.Print("starting server...")

	router.Run()
}
