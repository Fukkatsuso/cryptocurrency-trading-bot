package repository

import "github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"

type UserRepository interface {
	Save(user *model.User) error
	FindByID(id string) (*model.User, error)
}
