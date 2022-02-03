package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/external/bitflyer"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/usecase"
)

func Run() {
	// repository
	userRepository := persistence.NewUserRepository(config.DB)
	sessionRepository := persistence.NewSessionRepository(config.DB)
	candleRepository := persistence.NewCandleRepository(config.DB, config.CandleTableName, config.TimeFormat)
	signalEventRepository := persistence.NewSignalEventRepository(config.DB, config.TimeFormat)
	tradeParamsRepository := persistence.NewTradeParamsRepository(config.DB)
	cookie := persistence.NewCookie("cryptobot", "/", 60*30, config.SecureCookie)
	// repository (bitflyer)
	bitflyerClient := bitflyer.NewClient(config.APIKey, config.APISecret)
	balanceRepository := bitflyer.NewBitFlyerBalanceRepository(bitflyerClient)

	// service
	authService := service.NewAuthService(userRepository, sessionRepository)
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewMRBaseDataFrameService(indicatorService)

	// usecase
	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)
	tradeParamsUsecase := usecase.NewTradeParamsUsecase(tradeParamsRepository)
	balanceUsecase := usecase.NewBalanceUsecase(balanceRepository)

	// handler
	authHandler := handler.NewAuthHandler(cookie, authService)
	dataFrameHandler := handler.NewDataFrameHandler(dataFrameUsecase)
	tradeParamsHandler := handler.NewTradeParamsHandler(tradeParamsUsecase)
	balanceHandler := handler.NewBalanceHandler(balanceUsecase)

	http.HandleFunc("/api/login", authHandler.Login())
	http.HandleFunc("/api/logout", authHandler.Logout())
	http.HandleFunc("/api/candle", dataFrameHandler.Get(config.ProductCode))
	http.HandleFunc("/admin/api/trade-params", AuthGuardHandlerFunc(tradeParamsHandler.HandlerFunc(), authHandler))
	http.HandleFunc("/admin/api/balance", AuthGuardHandlerFunc(balanceHandler.Get(), authHandler))

	http.HandleFunc("/", PageHandlerFunc("view/index.html"))
	http.HandleFunc("/login", PageHandlerFunc("view/login.html"))
	http.HandleFunc("/admin", AuthGuardHandlerFunc(PageHandlerFunc("view/admin.html"), authHandler))
	http.Handle("/view/admin.html", http.RedirectHandler("/admin", http.StatusFound))
	http.Handle("/view/", http.StripPrefix("/view/", http.FileServer(http.Dir("view/"))))

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

func PageHandlerFunc(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(filepath))
		if err := tmpl.Execute(w, nil); err != nil {
			fmt.Println("[PageHandlerFunc]", err)
		}
	}
}

// ログイン済みユーザだけ受け付ける
func AuthGuardHandlerFunc(target http.HandlerFunc, ah handler.AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ah.LoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		target.ServeHTTP(w, r)
	}
}
