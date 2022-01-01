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

	sessionRepository := persistence.NewSessionRepository(db)

	// test user
	testUser := model.NewUser("test", "QWERTYUIOP", "")
	createTestUser(db, testUser)

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

func createTestUser(db persistence.DB, user *model.User) error {
	cmd := `
        INSERT INTO users
            (id, password_hash, session_id_hash)
        VALUES
            (?, ?, ?)
    `
	_, err := db.Exec(cmd, user.ID(), user.Password(), user.SessionID())
	return err
}
