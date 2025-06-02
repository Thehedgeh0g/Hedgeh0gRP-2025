package repository

import (
	"auth/pkg/infrastucture/protokey"
	"context"
	"encoding/json"
	"errors"

	"auth/pkg/app/model"
)

func NewUserProtoKeyRepository(pkdb *protokey.ProtoKeyClient) model.UserRepository {
	return &userRepository{
		protoKeyClient: pkdb,
		ctx:            context.Background(),
	}
}

type userSerializable struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userRepository struct {
	protoKeyClient *protokey.ProtoKeyClient
	ctx            context.Context
}

func (r *userRepository) FindByLogin(login string) (*model.User, error) {
	formattedData, err := r.protoKeyClient.Get(login)
	if err != nil {
		if errors.Is(err, protokey.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, err
	}
	userData := userSerializable{}
	err = json.Unmarshal([]byte(formattedData), &userData)
	user := model.LoadUser(
		userData.Login,
		userData.Password,
	)
	return &user, nil
}

func (r *userRepository) Store(user model.User) error {
	userData := userSerializable{
		Login:    user.Login(),
		Password: user.Password(),
	}
	formattedData, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return r.protoKeyClient.Set(user.Login(), string(formattedData))
}
