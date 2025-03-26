package model

import (
	"errors"
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
	FindByHash(hash string) (Text, error)
}

type Text struct {
	hash       string
	text       string
	similarity bool
	rank       *float64
}

func NewText(hash, text string) Text {
	return Text{
		hash:       hash,
		text:       text,
		similarity: false,
	}
}

func BuildTextFromSavedData(hash, text string, similarity bool, rank float64) Text {
	return Text{
		hash:       hash,
		text:       text,
		similarity: similarity,
		rank:       &rank,
	}
}

func (t *Text) GetRank() float64 {
	if t.rank != nil {
		return *t.rank
	}
	return 0
}

func (t *Text) GetHash() string {
	return t.hash
}

func (t *Text) GetText() string {
	return t.text
}

func (t *Text) GetSimilarity() bool {
	return t.similarity
}

func (t *Text) SetSimilarity(similarity bool) {
	t.similarity = similarity
}
