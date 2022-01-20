package repository

type SessionRepository interface {
	Save(userID string, sessionID string) error
	FindByUserID(userID string) (string, error)
	Delete(userID string) error
}
