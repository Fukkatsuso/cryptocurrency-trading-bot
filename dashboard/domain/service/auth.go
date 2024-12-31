package service

// type AuthService interface {
// 	Login(userID string, password string) (string, error)
// 	Logout(userID string) error
// 	LoggedIn(userID string, sessionID string) bool
// }

// type authService struct {
// 	userRepository    repository.UserRepository
// 	sessionRepository repository.SessionRepository
// }

// func NewAuthService(ur repository.UserRepository, sr repository.SessionRepository) AuthService {
// 	return &authService{
// 		userRepository:    ur,
// 		sessionRepository: sr,
// 	}
// }

// func (as *authService) Login(userID string, password string) (string, error) {
// 	user, err := as.userRepository.FindByID(userID)
// 	if err != nil {
// 		return "", err
// 	}

// 	if err := model.CompareHashAndPassword(user.Password(), password); err != nil {
// 		return "", errors.New("password is not correct")
// 	}

// 	sessionID, err := model.NewSessionID()
// 	if err != nil {
// 		return "", err
// 	}
// 	sessionIdHash, err := model.SessionIdHash(sessionID)
// 	if err != nil {
// 		return "", err
// 	}
// 	if err := as.sessionRepository.Save(userID, sessionIdHash); err != nil {
// 		return "", err
// 	}

// 	return sessionID, nil
// }

// func (as *authService) Logout(userID string) error {
// 	if err := as.sessionRepository.Delete(userID); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (as *authService) LoggedIn(userID string, sessionID string) bool {
// 	if sessionID == "" {
// 		return false
// 	}

// 	sessionIdHash, err := as.sessionRepository.FindByUserID(userID)
// 	if err != nil {
// 		return false
// 	}

// 	return model.CompareHashAndSessionID(sessionIdHash, sessionID) == nil
// }
