package jwt

import (
	"auth/pkg/app/model"
	"auth/pkg/app/service"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func NewTokenService(secret string, tokenRepository model.TokenRepository) service.TokenService {
	return &tokenService{
		secret:          secret,
		tokenRepository: tokenRepository,
	}
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type tokenService struct {
	secret          string
	tokenRepository model.TokenRepository
}

func (t *tokenService) CreateToken(login string, expirationTimeDur time.Duration) (string, time.Time, error) {
	expirationTime := time.Now().Add(expirationTimeDur)
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (t *tokenService) SaveToken(tokenString, login string) error {
	return t.tokenRepository.Store(tokenString, login)
}

func (t *tokenService) ParseToken(tokenString string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	if err != nil {
		return "", errors.New("token is invalid")
	}
	return claims.Login, nil
}

func (t *tokenService) DeleteToken(login string) error {
	return t.tokenRepository.Store("", login)
}
