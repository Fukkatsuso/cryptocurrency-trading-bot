package main

import (
	"fmt"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/router"
)

func main() {
	fmt.Println("starting server...")

	router.Run()
}
