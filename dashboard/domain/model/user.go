package model

type User struct {
	id        string
	password  string
	sessionID string
}

func NewUser(id string, password string, sessionID string) *User {
	if id == "" {
		return nil
	}

	if password == "" {
		return nil
	}

	return &User{
		id:        id,
		password:  password,
		sessionID: sessionID,
	}
}

func (user *User) Password() string {
	return user.password
}

func PasswordDigest(password string) string {
	return password
}

func NewSessionID() string {
	return "sessionID"
}

func SessionIdDigest(id string) string {
	return id
}
