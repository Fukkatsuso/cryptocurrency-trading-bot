package model

// type User struct {
// 	id        string
// 	password  string
// 	sessionID string
// }

// func NewUser(id string, password string, sessionID string) *User {
// 	if id == "" {
// 		return nil
// 	}

// 	if password == "" {
// 		return nil
// 	}

// 	return &User{
// 		id:        id,
// 		password:  password,
// 		sessionID: sessionID,
// 	}
// }

// func (user *User) ID() string {
// 	return user.id
// }

// func (user *User) Password() string {
// 	return user.password
// }

// func (user *User) SessionID() string {
// 	return user.sessionID
// }

// func PasswordHash(password string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hash), err
// }

// func CompareHashAndPassword(hash string, password string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// }

// func NewSessionID() (string, error) {
// 	u, err := uuid.NewRandom()
// 	if err != nil {
// 		return "", err
// 	}
// 	return u.String(), nil
// }

// func SessionIdHash(sessionID string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(sessionID), bcrypt.DefaultCost)
// 	return string(hash), err
// }

// func CompareHashAndSessionID(hash string, sessionID string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(sessionID))
// }
