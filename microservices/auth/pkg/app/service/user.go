package service

import (
	"auth/pkg/app/model"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
)

type UserService struct {
	userRepository model.UserRepository
}

func NewUserService(repository model.UserRepository) *UserService {
	return &UserService{userRepository: repository}
}

func (u *UserService) CreateUser(login, password string) error {
	user, err := u.userRepository.FindByLogin(login)
	if err != nil {
		return err
	}
	if user != nil {
		if user.Password() != hashString(password) {
			return errors.New("password incorrect")
		}
		return nil
	}

	return u.userRepository.Store(model.NewUser(login, hashString(password)))
}

func (a *UserService) Authenticate(login, password string) (bool, error) {
	user, err := a.userRepository.FindByLogin(login)
	if err != nil {
		return false, err
	}

	if user.Login() != login {
		log.Printf("User not exist %s", user.Login())
		return false, nil
	}

	if user.Password() != hashString(password) {
		log.Printf("Password not match %s", user.Password())
		return false, nil
	}

	return true, nil
}

func hashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
