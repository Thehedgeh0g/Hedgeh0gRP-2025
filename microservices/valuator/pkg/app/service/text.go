package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"

	"valuator/pkg/app/model"
)

type TextService struct {
	textRepository model.TextRepository
}

func NewTextService(textRepository model.TextRepository) *TextService {
	return &TextService{textRepository: textRepository}
}

func (s *TextService) ProcessText(data string) (uuid.UUID, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedStr := hex.EncodeToString(hash.Sum(nil))
	id, err := s.textRepository.FindByData(hashedStr)
	if err != nil && !errors.Is(err, model.ErrTextNotFound) {
		return uuid.MustParse(id), err
	}
	// DONE: хешировать данные, чтобы хранить тексты в ед. экземпляре.
	if errors.Is(err, model.ErrTextNotFound) {
		text := model.NewText(hashedStr)
		return text.GetID(), s.textRepository.Store(text)
	}

	text, err := s.textRepository.FindByID(uuid.MustParse(id))
	if err != nil && !errors.Is(err, model.ErrTextNotFound) {
		return uuid.MustParse(id), err
	}
	text.SetSimilarity(true)
	return text.GetID(), s.textRepository.Store(text)
}
