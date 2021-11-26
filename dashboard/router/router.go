package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/controller"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func Run() {
	http.HandleFunc("/", controller.IndexPageHandler)
	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("view/"))))

	// repository
	candleRepository := persistence.NewCandleRepository(config.DB, config.CandleTableName, config.TimeFormat)
	signalEventRepository := persistence.NewSignalEventRepository(config.DB, config.TimeFormat)

	// service
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)

	// usecase
	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)

	// handler
	dataFrameHandler := handler.NewDataFrameHandler(dataFrameUsecase)

	http.HandleFunc("/api/candle", dataFrameHandler.Get(config.ProductCode, 0.01))

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
