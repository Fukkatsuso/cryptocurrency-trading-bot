package model

type User struct {
	id       string
	password string
}

func NewUser(id string, password string) *User {
	if id == "" {
		return nil
	}

	if password == "" {
		return nil
	}

	return &User{
		id:       id,
		password: password,
	}
}
