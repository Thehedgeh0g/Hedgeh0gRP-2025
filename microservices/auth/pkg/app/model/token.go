package model

type TokenRepository interface {
	FindByLogin(login string) (string, error)
	Store(token, login string) error
}
