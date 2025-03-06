package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strconv"

	"valuator/pkg/app/model"
)

var (
	textRank       = "rank-"
	textSimilarity = "similarity-"
	textValue      = "value-"
)

type TextRepository struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewTextRepository(redisClient *redis.Client) *TextRepository {
	return &TextRepository{
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// Store сохраняет текст в Redis. Устанавливает ранг и схожесть для текста, если необходимо.
func (t *TextRepository) Store(text model.Text) error {
	if err := t.redisClient.Set(t.ctx, textValue+text.GetData(), text.GetID().String(), 0).Err(); err != nil {
		return err
	}
	if err := t.redisClient.Set(t.ctx, textRank+text.GetID().String(), text.GetRank(), 0).Err(); err != nil {
		return err
	}
	if err := t.redisClient.Set(t.ctx, textSimilarity+text.GetID().String(), text.GetSimilarity(), 0).Err(); err != nil {
		return err
	}
	if !text.GetSimilarity() {
		return t.redisClient.Set(t.ctx, text.GetData(), text.GetID().String(), 0).Err()
	}
	return nil
}

// FindByData ищет текст в Redis по его данным.
// DONE: переименовать
func (t *TextRepository) FindByHash(data string) (string, error) {
	id, err := t.redisClient.Get(t.ctx, textValue+data).Result()
	if errors.Is(err, redis.Nil) {
		return "", model.ErrTextNotFound
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *TextRepository) FindByID(id uuid.UUID) (model.Text, error) {
	rankStr, err := t.redisClient.Get(t.ctx, textRank+id.String()).Result()
	if err != nil {
		return model.Text{}, err
	}
	similarityStr, err := t.redisClient.Get(t.ctx, textSimilarity+id.String()).Result()
	if err != nil {
		return model.Text{}, err
	}
	rank, err := strconv.ParseFloat(rankStr, 64)
	if err != nil {
		return model.Text{}, err
	}
	similarity, err := strconv.ParseBool(similarityStr)
	if err != nil {
		return model.Text{}, err
	}
	return model.BuildTextFromSavedData("", id, similarity, rank), nil
}
