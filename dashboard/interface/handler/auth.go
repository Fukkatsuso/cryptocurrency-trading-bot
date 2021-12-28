package handler

import (
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/service"
	"github.com/gorilla/securecookie"
)

type AuthHandler interface {
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
}

type authHandler struct {
	cookieName   string
	cookiePath   string
	cookieMaxAge int // ブラウザでcookieが削除されるまでの秒数
	secureCookie *securecookie.SecureCookie
	authService  service.AuthService
}

func NewAuthHandler(cookieName string, cookiePath string, cookieMaxAge int, secureCookie *securecookie.SecureCookie, as service.AuthService) AuthHandler {
	return &authHandler{
		cookieName:   cookieName,
		cookiePath:   cookiePath,
		cookieMaxAge: cookieMaxAge,
		secureCookie: secureCookie,
		authService:  as,
	}
}

func (ah *authHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
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
		encoded, err := ah.secureCookie.Encode(ah.cookieName, cookieValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:     ah.cookieName,
			Value:    encoded,
			Path:     ah.cookiePath,
			MaxAge:   ah.cookieMaxAge,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
	}
}

func (ah *authHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// cookieの値を取得
		cookie, err := r.Cookie(ah.cookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		value := make(map[string]string)
		if err = ah.secureCookie.Decode(ah.cookieName, cookie.Value, &value); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID := value["userID"]
		if err := ah.authService.Logout(userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// cookieを削除
		cookie = &http.Cookie{
			Name:     ah.cookieName,
			Path:     ah.cookiePath,
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusOK)
	}
}
