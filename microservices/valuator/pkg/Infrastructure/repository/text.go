package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"valuator/pkg/app/model"
)

type dbText struct {
	Rank       float64 `json:"rank"`
	Text       string  `json:"text"`
	Similarity bool    `json:"similarity"`
}

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
	textData := dbText{
		Rank:       text.GetRank(),
		Text:       text.GetText(),
		Similarity: text.GetSimilarity(),
	}
	fmt.Println(text.GetText())
	formattedData, _ := json.Marshal(textData)
	fmt.Println(string(formattedData))

	return t.redisClient.Set(t.ctx, text.GetHash(), formattedData, 0).Err()
}

func (t *TextRepository) FindByHash(hash string) (model.Text, error) {
	formattedData, err := t.redisClient.Get(t.ctx, hash).Result()
	if errors.Is(err, redis.Nil) {
		return model.Text{}, model.ErrTextNotFound
	}
	if err != nil {
		return model.Text{}, err
	}
	textData := dbText{}
	err = json.Unmarshal([]byte(formattedData), &textData)
	if err != nil {
		return model.Text{}, err
	}

	return model.BuildTextFromSavedData(hash, textData.Text, textData.Similarity, textData.Rank), nil
}
