package persistence

import (
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
	"github.com/gorilla/securecookie"
)

type cookie struct {
	cookieName   string
	cookiePath   string
	cookieMaxAge int // ブラウザでcookieが削除されるまでの秒数
	secureCookie *securecookie.SecureCookie
}

func NewCookie(cookieName string, cookiePath string, cookieMaxAge int, secureCookie *securecookie.SecureCookie) repository.Cookie {
	return &cookie{
		cookieName:   cookieName,
		cookiePath:   cookiePath,
		cookieMaxAge: cookieMaxAge,
		secureCookie: secureCookie,
	}
}

func (c *cookie) Set(w http.ResponseWriter, value map[string]string) error {
	encoded, err := c.secureCookie.Encode(c.cookieName, value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     c.cookieName,
		Value:    encoded,
		Path:     c.cookiePath,
		MaxAge:   c.cookieMaxAge,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}

func (c *cookie) GetValue(r *http.Request) (map[string]string, error) {
	cookie, err := r.Cookie(c.cookieName)
	if err != nil {
		return nil, err
	}

	value := make(map[string]string)
	err = c.secureCookie.Decode(c.cookieName, cookie.Value, &value)

	return value, err
}

func (c *cookie) Delete(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     c.cookieName,
		Path:     c.cookiePath,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
