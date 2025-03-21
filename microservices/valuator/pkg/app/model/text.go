package model

import (
	"errors"
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
	re := regexp.MustCompile(`[A-Za-zА-Яа-я]`)
	alphabetCount := float64(len(re.FindAllString(t.text, -1)))
	totalCount := float64(len(t.text))

	return 1 - (alphabetCount / totalCount)
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
