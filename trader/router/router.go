package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/external/slack"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/trader/usecase"
)

func Run() {
	// repository
	candleRepository := persistence.NewCandleRepository(config.DB, config.CandleTableName, config.TimeFormat)
	signalEventRepository := persistence.NewSignalEventRepository(config.DB, config.TimeFormat)
	tradeParamsRepository := persistence.NewTradeParamsRepository(config.DB)
	// repository (bitflyer)
	bitflyerClient := bitflyer.NewClient(config.APIKey, config.APISecret)
	tickerRepository := bitflyer.NewBitflyerTickerRepository(bitflyerClient)
	balanceRepository := bitflyer.NewBitFlyerBalanceRepository(bitflyerClient)
	orderRepository := bitflyer.NewBitflyerOrderRepository(bitflyerClient)
	// repository (slack)
	slackClient := slack.NewClient(config.SlackBotToken, config.SlackChannelID)
	notificationRepository := slack.NewSlackNotificationRepository(slackClient, config.LocalTime)

	// service
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)
	tradeParamsService := service.NewTradeParamsService(tradeParamsRepository, dataFrameService)
	tradeService := service.NewTradeService(balanceRepository, tickerRepository, orderRepository, signalEventRepository, candleService, dataFrameService, tradeParamsService)
	notificationService := service.NewNotificationService(notificationRepository)

	// usecase
	candleUsecase := usecase.NewCandleUsecase(candleService, tickerRepository)
	tradeUsecase := usecase.NewTradeUsecase(signalEventService, tradeService, notificationService)

	// handler
	candleHandler := handler.NewCandleHandler(candleUsecase)
	tradeHandler := handler.NewTradeHandler(tradeUsecase)

	http.HandleFunc("/fetch-ticker", candleHandler.UpdateCandle(config.ProductCode))
	http.HandleFunc("/trade", tradeHandler.Trade(config.ProductCode, 365))

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
