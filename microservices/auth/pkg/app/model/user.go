package model

type User struct {
	login    string
	password string
}

type UserRepository interface {
	FindByLogin(login string) (*User, error)
	Store(user User) error
}

func NewUser(login string, password string) User {
	return User{
		login:    login,
		password: password,
	}
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Password() string {
	return u.password
}

func LoadUser(login string, password string) User {
	return User{
		login:    login,
		password: password,
	}
}
