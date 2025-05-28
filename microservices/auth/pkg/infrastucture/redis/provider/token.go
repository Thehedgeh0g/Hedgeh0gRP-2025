package provider

import (
	"auth/pkg/app/provider"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

func NewTokenProvider(rdb *redis.Client) provider.TokenProvider {
	return &tokenProvider{
		redisClient: rdb,
		ctx:         context.Background(),
	}
}

type tokenSerializable struct {
	TokenString string `json:"token"`
}

type tokenProvider struct {
	redisClient *redis.Client
	ctx         context.Context
}

func (t *tokenProvider) GetTokenByLogin(login string) (string, error) {
	tokenKey := "Token:" + login
	formattedData, err := t.redisClient.Get(t.ctx, tokenKey).Result()
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
