package service

import (
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type UserService interface {
	Login(id string, password string) (string, error)
	Logout(id string) error
	LoggedIn(id string, password string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) Login(id string, password string) (string, error) {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return "", err
	}

	passwordDigest := model.PasswordDigest(password)
	if passwordDigest != user.Password() {
		return "", errors.New("password is not correct")
	}

	sessionID := model.NewSessionID()

	user = model.NewUser(id, passwordDigest, sessionID)
	if err := us.userRepository.Save(user); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (us *userService) Logout(id string) error {
	return nil
}

func (us *userService) LoggedIn(id string, password string) bool {
	return false
}
