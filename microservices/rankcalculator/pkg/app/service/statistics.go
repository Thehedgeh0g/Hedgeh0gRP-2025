package service

import (
	"rankcalculator/pkg/app/dispatcher"
	appevent "rankcalculator/pkg/app/event"
	"regexp"

	"rankcalculator/pkg/app/model"
)

type StatisticsService interface {
	CalculateRank(hash string) error
}

func NewStatisticsService(
	textRepo model.TextRepository,
	eventDispatcher dispatcher.EventDispatcher,
) StatisticsService {
	return &statisticsService{
		textRepo:        textRepo,
		eventDispatcher: eventDispatcher,
	}
}

type statisticsService struct {
	textRepo        model.TextRepository
	eventDispatcher dispatcher.EventDispatcher
}

func (service *statisticsService) CalculateRank(hash string) error {
	text, err := service.textRepo.FindByHash(hash)
	if err != nil {
		return err
	}
	textData := text.GetText()
	re := regexp.MustCompile(`[A-Za-zА-Яа-я]`)
	alphabetCount := float64(len(re.FindAllString(textData, -1)))
	totalCount := float64(len(textData))

	rank := 1 - (alphabetCount / totalCount)
	text.SetRank(rank)
	err = service.textRepo.Store(text)
	if err != nil {
		return err
	}
	return service.eventDispatcher.Dispatch(appevent.NewRankCalculatedEvent(text.GetHash(), rank))
}
