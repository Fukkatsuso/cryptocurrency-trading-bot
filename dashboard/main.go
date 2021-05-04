package main

import (
	"fmt"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/router"
)

func main() {
	fmt.Println("starting server...")

	router.Run()
}
