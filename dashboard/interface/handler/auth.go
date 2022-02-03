package handler

import (
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
)

type AuthHandler interface {
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	LoggedIn(r *http.Request) bool
}

type authHandler struct {
	cookie      repository.Cookie
	authService service.AuthService
}

func NewAuthHandler(cookie repository.Cookie, as service.AuthService) AuthHandler {
	return &authHandler{
		cookie:      cookie,
		authService: as,
	}
}

func (ah *authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "this method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.FormValue("userId")
		password := r.FormValue("password")

		sessionID, err := ah.authService.Login(userID, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// userIDとsessionIDをCookieにセット
		cookieValue := map[string]string{
			"userID":    userID,
			"sessionID": sessionID,
		}
		err = ah.cookie.Set(w, cookieValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}
}

func (ah *authHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "this method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		// cookieの値を取得
		cookieValue, err := ah.cookie.GetValue(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := cookieValue["userID"]
		if err := ah.authService.Logout(userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// cookieを削除
		ah.cookie.Delete(w)

		http.Redirect(w, r, "/", http.StatusOK)
	}
}

func (ah *authHandler) LoggedIn(r *http.Request) bool {
	// cookieの値を取得
	cookieValue, err := ah.cookie.GetValue(r)
	if err != nil {
		return false
	}

	userID := cookieValue["userID"]
	sessionID := cookieValue["sessionID"]

	return ah.authService.LoggedIn(userID, sessionID)
}
