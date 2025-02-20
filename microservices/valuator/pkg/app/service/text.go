package service

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"valuator/pkg/app/model"
)

type TextService struct {
	textRepository model.TextRepository
}

func NewTextService(textRepository model.TextRepository) *TextService {
	return &TextService{textRepository: textRepository}
}

func (s *TextService) ProcessText(data string) (uuid.UUID, error) {
	text := model.NewText(strings.ReplaceAll(data, "\n", "\\n"))

	_, err := s.textRepository.FindByData(text)
	if err != nil && !errors.Is(err, model.ErrTextNotFound) {
		return text.GetID(), err
	}

	if errors.Is(err, model.ErrTextNotFound) {
		return text.GetID(), s.textRepository.Store(text)
	}

	text.SetSimilarity(true)
	return text.GetID(), s.textRepository.Store(text)
}
