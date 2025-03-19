package model

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

func BuildTextFromSavedData(hash, text string, similarity bool, rank float64) Text {
	return Text{
		hash:       hash,
		text:       text,
		similarity: similarity,
		rank:       &rank,
	}
}

func (t *Text) GetRank() float64 {
	return *t.rank
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

func (t *Text) SetRank(newRank float64) {
	t.rank = &newRank
}
