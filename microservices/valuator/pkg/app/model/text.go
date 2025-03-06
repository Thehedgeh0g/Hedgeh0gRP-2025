package model

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
)

var (
	ErrTextNotFound = errors.New("text not found")
)

type TextStats struct {
	Rank       float64
	Similarity bool
}

type TextRepository interface {
	Store(text Text) error
	FindByHash(data string) (string, error)
	FindByID(id uuid.UUID) (Text, error)
}

type Text struct {
	id         uuid.UUID
	data       string
	similarity bool
	rank       *float64
}

func NewText(data string) Text {
	return Text{
		id:         uuid.New(),
		data:       data,
		similarity: false,
	}
}

func BuildTextFromSavedData(data string, id uuid.UUID, similarity bool, rank float64) Text {
	return Text{
		id:         id,
		data:       data,
		similarity: similarity,
		rank:       &rank,
	}
}

func (t *Text) GetRank() float64 {
	if t.rank != nil {
		return *t.rank
	}
	re := regexp.MustCompile(`[A-Za-zА-Яа-я]`)
	alphabetCount := float64(len(re.FindAllString(t.data, -1)))
	totalCount := float64(len(t.data))

	return 1 - (alphabetCount / totalCount)
}

func (t *Text) GetData() string {
	return t.data
}

func (t *Text) GetID() uuid.UUID {
	return t.id
}

func (t *Text) GetSimilarity() bool {
	return t.similarity
}

func (t *Text) SetSimilarity(similarity bool) {
	t.similarity = similarity
}
