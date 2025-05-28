package service

import (
	"time"
)

type TokenService interface {
	CreateToken(login string, expirationTimeDur time.Duration) (string, time.Time, error)
	SaveToken(tokenString, login string) error
	ParseToken(tokenString string) (string, error)
	DeleteToken(login string) error
}
