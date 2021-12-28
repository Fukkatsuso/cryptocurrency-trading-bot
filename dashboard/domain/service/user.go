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
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewUserService(ur repository.UserRepository, sr repository.SessionRepository) UserService {
	return &userService{
		userRepository:    ur,
		sessionRepository: sr,
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
	if err := us.sessionRepository.Save(id, sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (us *userService) Logout(id string) error {
	if err := us.sessionRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (us *userService) LoggedIn(id string, sessionID string) bool {
	if sessionID == "" {
		return false
	}

	sessionIdDigest, err := us.sessionRepository.FindByUserID(id)
	if err != nil {
		return false
	}

	return model.SessionIdDigest(sessionID) == sessionIdDigest
}
