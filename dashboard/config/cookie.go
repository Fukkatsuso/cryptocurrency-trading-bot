package config

import (
	"fmt"
	"os"

	"github.com/gorilla/securecookie"
)

var SecureCookie *securecookie.SecureCookie

var (
	COOKIE_HASHKEY  = []byte(os.Getenv("COOKIE_HASHKEY"))  // at least 32 bytes long
	COOKIE_BLOCKKEY = []byte(os.Getenv("COOKIE_BLOCKKEY")) // 16 bytes (AES-128) or 32 bytes (AES-256) long
)

const CookieName = "fukkatsuso-cryptocurrency"

func init() {
	if len(COOKIE_HASHKEY) < 32 {
		fmt.Println("COOKIE_HASHKEY should be at least 32 bytes long")
	}

	if len(COOKIE_BLOCKKEY) != 16 && len(COOKIE_BLOCKKEY) != 32 {
		fmt.Println("COOKIE_BLOCKKEY should be 16 bytes (AES-128) or 32 bytes (AES-256) long")
	}

	SecureCookie = securecookie.New(COOKIE_HASHKEY, COOKIE_BLOCKKEY)
}
