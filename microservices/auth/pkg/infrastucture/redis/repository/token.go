package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"auth/pkg/app/model"
)

func NewTokenRepository(rdb *redis.Client) model.TokenRepository {
	return &tokenRepository{
		redisClient: rdb,
		ctx:         context.Background(),
	}
}

type tokenSerializable struct {
	TokenString string `json:"token"`
}

type tokenRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (r *tokenRepository) FindByLogin(login string) (string, error) {
	tokenKey := "Token:" + login
	formattedData, err := r.redisClient.Get(r.ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	token := tokenSerializable{}
	err = json.Unmarshal([]byte(formattedData), &token)

	return token.TokenString, nil
}

func (r *tokenRepository) Store(token, login string) error {
	tokenData := tokenSerializable{
		token,
	}
	formattedData, err := json.Marshal(tokenData)
	if err != nil {
		return err
	}
	tokenKey := "Token:" + login
	return r.redisClient.Set(r.ctx, tokenKey, formattedData, 0).Err()
}
