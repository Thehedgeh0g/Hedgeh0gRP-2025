package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"rankcalculator/pkg/app/model"
)

type dbText struct {
	Rank       float64 `json:"rank"`
	Text       string  `json:"text"`
	Similarity bool    `json:"similarity"`
}

type TextRepository struct {
	shardManager *ShardManager
	ctx          context.Context
}

func NewTextRepository(shardManager *ShardManager) *TextRepository {
	return &TextRepository{
		shardManager: shardManager,
		ctx:          context.Background(),
	}
}

func (t *TextRepository) Store(text model.Text) error {
	client, region, isNewRegion, err := t.shardManager.GetClientByHash(text.GetHash())
	if err != nil {
		return err
	}
	if client == nil {
		return model.ErrTextNotFound
	}
	fmt.Println(client.String())

	textData := dbText{
		Rank:       text.GetRank(),
		Text:       text.GetText(),
		Similarity: text.GetSimilarity(),
	}

	formattedData, _ := json.Marshal(textData)
	if err = client.Set(t.ctx, text.GetHash(), formattedData, 0).Err(); err != nil {
		return err
	}

	if isNewRegion {
		return t.shardManager.SetRegionForHash(text.GetHash(), region)
	}
	return nil
}

func (t *TextRepository) FindByHash(hash string) (model.Text, error) {
	client, _, _, err := t.shardManager.GetClientByHash(hash)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Text{}, model.ErrTextNotFound
		}
		return model.Text{}, err
	}
	if client == nil {
		return model.Text{}, model.ErrTextNotFound
	}

	formattedData, err := client.Get(t.ctx, hash).Result()
	if errors.Is(err, redis.Nil) {
		return model.Text{}, model.ErrTextNotFound
	}
	if err != nil {
		return model.Text{}, err
	}

	var textData dbText
	if err := json.Unmarshal([]byte(formattedData), &textData); err != nil {
		return model.Text{}, err
	}

	return model.BuildTextFromSavedData(hash, textData.Text, textData.Similarity, textData.Rank), nil
}
