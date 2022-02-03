package persistence

import (
	"net/http"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
	"github.com/gorilla/securecookie"
)

type cookie struct {
	name         string
	path         string
	maxAge       int // ブラウザでcookieが削除されるまでの秒数
	secureCookie *securecookie.SecureCookie
}

func NewCookie(name string, path string, maxAge int, secureCookie *securecookie.SecureCookie) repository.Cookie {
	return &cookie{
		name:         name,
		path:         path,
		maxAge:       maxAge,
		secureCookie: secureCookie,
	}
}

func (c *cookie) Set(w http.ResponseWriter, value map[string]string) error {
	encoded, err := c.secureCookie.Encode(c.name, value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     c.name,
		Value:    encoded,
		Path:     c.path,
		MaxAge:   c.maxAge,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	return nil
}

func (c *cookie) GetValue(r *http.Request) (map[string]string, error) {
	cookie, err := r.Cookie(c.name)
	if err != nil {
		return nil, err
	}

	value := make(map[string]string)
	err = c.secureCookie.Decode(c.name, cookie.Value, &value)

	return value, err
}

func (c *cookie) Delete(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     c.name,
		Path:     c.path,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
