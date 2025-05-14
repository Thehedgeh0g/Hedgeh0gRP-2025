package model

import "errors"

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

type Text interface {
	GetRank() float64
	GetHash() string
	GetText() string
	GetSimilarity() bool
	SetSimilarity(similarity bool)
	SetRank(newRank float64)
}

type text struct {
	hash       string
	text       string
	similarity bool
	rank       *float64
}

func BuildTextFromSavedData(hash, textData string, similarity bool, rank float64) Text {
	return &text{
		hash:       hash,
		text:       textData,
		similarity: similarity,
		rank:       &rank,
	}
}

func (t *text) GetRank() float64 {
	return *t.rank
}

func (t *text) GetHash() string {
	return t.hash
}

func (t *text) GetText() string {
	return t.text
}

func (t *text) GetSimilarity() bool {
	return t.similarity
}

func (t *text) SetSimilarity(similarity bool) {
	t.similarity = similarity
}

func (t *text) SetRank(newRank float64) {
	t.rank = &newRank
}
