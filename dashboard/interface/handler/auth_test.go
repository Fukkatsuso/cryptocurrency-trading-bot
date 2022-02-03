package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/interface/handler"
)

func TestAuth(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	userRepository := persistence.NewUserRepository(tx)
	sessionRepository := persistence.NewSessionRepository(tx)
	cookie := persistence.NewCookie("cryptobot", "/", 60*30, config.SecureCookie)

	authService := service.NewAuthService(userRepository, sessionRepository)

	authHandler := handler.NewAuthHandler(cookie, authService)

	// create testUser
	passwordHash, err := model.PasswordHash("password")
	if err != nil {
		t.Fatal(err.Error())
	}
	testUser := model.NewUser("test", passwordHash, "")
	if err := userRepository.Save(testUser); err != nil {
		t.Fatal(err.Error())
	}

	var cookies []*http.Cookie

	t.Run("not logged in", func(t *testing.T) {
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		if authHandler.LoggedIn(req) {
			t.Fatal("logged in")
		}
	})

	t.Run("login", func(t *testing.T) {
		ts := httptest.NewServer(authHandler.Login())
		defer ts.Close()

		form := url.Values{
			"userId":   {"test"},
			"password": {"password"},
		}

		resp, err := http.PostForm(ts.URL, form)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("resp.StatusCode != http.StatusOK")
		}

		cookies = resp.Cookies()

		// check
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		for _, c := range cookies {
			req.AddCookie(c)
		}
		if !authHandler.LoggedIn(req) {
			t.Fatal("not logged in")
		}
	})

	t.Run("logout", func(t *testing.T) {
		ts := httptest.NewServer(authHandler.Logout())
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL, nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		for _, c := range cookies {
			req.AddCookie(c)
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("resp.StatusCode != http.StatusOK")
		}

		cookies = resp.Cookies()

		// check
		req, err = http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err.Error())
		}
		for _, c := range cookies {
			req.AddCookie(c)
		}
		if authHandler.LoggedIn(req) {
			t.Fatal("logged in")
		}
	})
}
