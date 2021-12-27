package repository

type SessionRepository interface {
	Save(userID string, sessionID string) error
	Delete(userID string) error
}
