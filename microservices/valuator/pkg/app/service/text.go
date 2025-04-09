package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	appevent "valuator/pkg/app/event"
	"valuator/pkg/app/model"
)

type TextService struct {
	textRepository  model.TextRepository
	eventDispatcher appevent.EventDispatcher
}

func NewTextService(textRepository model.TextRepository, eventDispatcher appevent.EventDispatcher) *TextService {
	return &TextService{textRepository: textRepository, eventDispatcher: eventDispatcher}
}

func (s *TextService) EvaluateText(data string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedStr := hex.EncodeToString(hash.Sum(nil))
	text, err := s.textRepository.FindByHash(hashedStr)
	if err != nil && !errors.Is(err, model.ErrTextNotFound) {
		return "", err
	}
	if !errors.Is(err, model.ErrTextNotFound) {
		text.SetSimilarity(true)
		err = s.textRepository.Store(text)
		if err != nil {
			return "", err
		}
	} else {
		text = model.NewText(hashedStr, data)
		err = s.textRepository.Store(text)
		if err != nil {
			return "", err
		}
		err = s.eventDispatcher.Dispatch(appevent.NewTextAddedEvent(hashedStr))
		err = s.eventDispatcher.Dispatch(appevent.NewSimilarityCalculatedEvent(hashedStr, text.GetSimilarity()))
	}
	return text.GetHash(), nil
}
