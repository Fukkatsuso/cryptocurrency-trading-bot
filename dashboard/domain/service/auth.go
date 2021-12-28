package service

import (
	"errors"

	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/model"
	"github.com/Fukkatsuso/cryptocurrency-trading-bot/dashboard/domain/repository"
)

type AuthService interface {
	Login(userID string, password string) (string, error)
	Logout(userID string) error
	LoggedIn(userID string, password string) bool
}

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewAuthService(ur repository.UserRepository, sr repository.SessionRepository) AuthService {
	return &authService{
		userRepository:    ur,
		sessionRepository: sr,
	}
}

func (as *authService) Login(userID string, password string) (string, error) {
	user, err := as.userRepository.FindByID(userID)
	if err != nil {
		return "", err
	}

	passwordDigest := model.PasswordDigest(password)
	if passwordDigest != user.Password() {
		return "", errors.New("password is not correct")
	}

	sessionID := model.NewSessionID()
	sessionIdDigest := model.SessionIdDigest(sessionID)
	if err := as.sessionRepository.Save(userID, sessionIdDigest); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (as *authService) Logout(userID string) error {
	if err := as.sessionRepository.Delete(userID); err != nil {
		return err
	}

	return nil
}

func (as *authService) LoggedIn(userID string, sessionID string) bool {
	if sessionID == "" {
		return false
	}

	sessionIdDigest, err := as.sessionRepository.FindByUserID(userID)
	if err != nil {
		return false
	}

	return model.SessionIdDigest(sessionID) == sessionIdDigest
}
