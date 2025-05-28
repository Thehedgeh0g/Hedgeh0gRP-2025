package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"auth/pkg/app/model"
)

func NewUserRepository(rdb *redis.Client) model.UserRepository {
	return &userRepository{
		redisClient: rdb,
		ctx:         context.Background(),
	}
}

type userSerializable struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type userRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (r *userRepository) FindByLogin(login string) (*model.User, error) {
	formattedData, err := r.redisClient.Get(r.ctx, login).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
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
	return r.redisClient.Set(r.ctx, user.Login(), formattedData, 0).Err()
}
