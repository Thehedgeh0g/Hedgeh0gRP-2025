package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ShardManager struct {
	ctx          context.Context
	region       string
	redisClient  *redis.Client
	dbClientsMap *map[string]*redis.Client
}

func NewShardManager(redisClient *redis.Client, dbClientsMap *map[string]*redis.Client, region string) *ShardManager {
	for _, client := range *dbClientsMap {
		fmt.Println("NewShardManager", client.String())
	}
	return &ShardManager{
		ctx:          context.Background(),
		region:       region,
		redisClient:  redisClient,
		dbClientsMap: dbClientsMap,
	}
}

// GetClientByHash возвращает клиента Redis по хэшу и флаг, была ли новая запись.
func (s *ShardManager) GetClientByHash(hash string) (*redis.Client, string, bool, error) {
	region, err := s.redisClient.Get(s.ctx, hash).Result()
	if errors.Is(err, redis.Nil) {
		region = s.region
		if region == "" {
			return nil, "", false, errors.New("region not found")
		}
		client, ok := (*s.dbClientsMap)[region] // Правильное использование
		if !ok {
			return nil, "", false, errors.New("unknown region: " + region)
		}
		fmt.Println(fmt.Sprintf("1 LOOKUP: %s; %s", region, client.String()))
		return client, region, true, nil
	}
	if err != nil {
		return nil, "", false, err
	}
	client, ok := (*s.dbClientsMap)[region] // Правильное использование
	if !ok {
		return nil, "", false, errors.New("unknown region: " + region)
	}
	fmt.Println(fmt.Sprintf("2 LOOKUP: %s; %s", region, client.String()))
	return client, region, false, nil
}

func (s *ShardManager) SetRegionForHash(hash, region string) error {
	fmt.Println("SET REGION HASH:", hash, "FOR REGION: ", region)
	return s.redisClient.Set(s.ctx, hash, region, 0).Err()
}
