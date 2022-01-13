package model_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestUser(t *testing.T) {
	tx := persistence.NewMySQLTransaction(config.DSN())
	defer tx.Rollback()

	t.Run("NewUser", func(t *testing.T) {
		var user *model.User

		user = model.NewUser("id", "password", "sessionID")
		if user == nil {
			t.Fatal("model.NewUser() returns nil")
		}

		user = model.NewUser("", "", "")
		if user != nil {
			t.Fatal("model.NewUser() returns not nil")
		}

		user = model.NewUser("", "password", "")
		if user != nil {
			t.Fatal("model.NewUser() returns not nil")
		}

		user = model.NewUser("id", "", "")
		if user != nil {
			t.Fatal("model.NewUser() returns not nil")
		}
	})

	t.Run("password hash", func(t *testing.T) {
		password := "password"
		passwordHash, err := model.PasswordHash(password)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = model.CompareHashAndPassword(passwordHash, password)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("sessionID hash", func(t *testing.T) {
		sessID, err := model.NewSessionID()
		if err != nil {
			t.Fatal(err.Error())
		}

		sessIdHash, err := model.SessionIdHash(sessID)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = model.CompareHashAndSessionID(sessIdHash, sessID)
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}
