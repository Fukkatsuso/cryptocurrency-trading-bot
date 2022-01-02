package persistence

import (
	"errors"
	"fmt"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type userRepository struct {
	db DB
}

func NewUserRepository(db DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Save(user *model.User) error {
	cmd := `
        INSERT INTO users
            (id, password_hash)
        VALUES
            (?, ?)
        ON DUPLICATE KEY UPDATE
            password_hash = VALUES(password_hash)
    `
	_, err := ur.db.Exec(cmd, user.ID(), user.Password())
	return err
}

func (ur *userRepository) FindByID(id string) (*model.User, error) {
	cmd := `
        SELECT
            id, password_hash, session_id_hash
        FROM
            users
        WHERE
            id = ?
    `
	row := ur.db.QueryRow(cmd, id)

	var userId, passwordHash, sessionIdHash string
	err := row.Scan(&userId, &passwordHash, &sessionIdHash)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(userId, passwordHash, sessionIdHash)
	if user == nil {
		return nil, errors.New(fmt.Sprint("invalid user:", userId, passwordHash, sessionIdHash))
	}
	return user, nil
}
