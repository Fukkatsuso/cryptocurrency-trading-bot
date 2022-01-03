package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
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

	// service
	authService := service.NewAuthService(userRepository, sessionRepository)
	candleService := service.NewCandleServicePerDay(config.LocalTime, config.TradeHour, candleRepository)
	signalEventService := service.NewSignalEventService(signalEventRepository)
	indicatorService := service.NewIndicatorService()
	dataFrameService := service.NewDataFrameService(indicatorService)

	// usecase
	dataFrameUsecase := usecase.NewDataFrameUsecase(candleService, signalEventService, dataFrameService)

	// handler
	authHandler := handler.NewAuthHandler("cryptobot", "/", 60*30, config.SecureCookie, authService)
	dataFrameHandler := handler.NewDataFrameHandler(dataFrameUsecase)

	http.HandleFunc("/api/login", authHandler.Login())
	http.HandleFunc("/api/logout", authHandler.Logout())
	http.HandleFunc("/api/candle", dataFrameHandler.Get(config.ProductCode, 0.01))

	http.HandleFunc("/", IndexPageHandler)
	http.HandleFunc("/login", LoginPageHandler)
	http.HandleFunc("/admin", AdminPageHandler(authHandler))
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

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[IndexPageHandler]", err)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view/login.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("[LoginPageHandler]", err)
	}
}

func AdminPageHandler(ah handler.AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ah.LoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		tmpl := template.Must(template.ParseFiles("view/admin.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			fmt.Println("[AdminPageHandler]", err)
		}
	}
}
