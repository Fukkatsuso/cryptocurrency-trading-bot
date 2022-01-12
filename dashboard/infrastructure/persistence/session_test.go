package persistence_test

import (
	"testing"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/config"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/infrastructure/persistence"
)

func TestSession(t *testing.T) {
	db := persistence.NewMySQLTransaction(config.DSN())
	defer db.Rollback()

	userRepository := persistence.NewUserRepository(db)
	sessionRepository := persistence.NewSessionRepository(db)

	// create test user
	testUser := model.NewUser("test", "QWERTYUIOP", "")
	userRepository.Save(testUser)

	sessionID := "abcdefghijklmnopqrstuvwxyz"

	t.Run("save", func(t *testing.T) {
		err := sessionRepository.Save(testUser.ID(), sessionID)
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("find by user id", func(t *testing.T) {
		sessionIdFound, err := sessionRepository.FindByUserID(testUser.ID())
		if err != nil {
			t.Fatal(err.Error())
		}
		if sessionIdFound != sessionID {
			t.Fatalf("%s != %s", sessionIdFound, sessionID)
		}
	})

	t.Run("delete", func(t *testing.T) {
		err := sessionRepository.Delete(testUser.ID())
		if err != nil {
			t.Fatal(err.Error())
		}

		sessionIdFound, err := sessionRepository.FindByUserID(testUser.ID())
		if err != nil {
			t.Fatal(err.Error())
		}
		if sessionIdFound != "" {
			t.Fatal("sessionID is not deleted")
		}
	})
}
