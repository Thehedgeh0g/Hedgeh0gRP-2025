package service

import (
	"fmt"
	"rankcalculator/pkg/app/dispatcher"
	appevent "rankcalculator/pkg/app/event"
	"rankcalculator/pkg/app/model"
	"rankcalculator/pkg/infrastructure/centrifugo"
	"regexp"
	"unicode/utf8"
)

type StatisticsService interface {
	CalculateRank(hash string) error
}

func NewStatisticsService(
	textRepo model.TextRepository,
	eventDispatcher dispatcher.EventDispatcher,
	centrifugoClient centrifugo.CentrifugoClient,
) StatisticsService {
	return &statisticsService{
		textRepo:         textRepo,
		eventDispatcher:  eventDispatcher,
		centrifugoClient: centrifugoClient,
	}
}

type statisticsService struct {
	textRepo         model.TextRepository
	eventDispatcher  dispatcher.EventDispatcher
	centrifugoClient centrifugo.CentrifugoClient
}

func (service *statisticsService) CalculateRank(hash string) error {
	text, err := service.textRepo.FindByHash(hash)
	if err != nil {
		return err
	}
	textData := text.GetText()
	re := regexp.MustCompile(`[A-Za-zА-Яа-яЁё]`)
	alphabetCount := float64(len(re.FindAllString(textData, -1)))
	fmt.Println(alphabetCount)
	totalCount := float64(utf8.RuneCountInString(textData))
	fmt.Println(totalCount)
	rank := 0.0
	if totalCount != 0 {
		rank = alphabetCount / totalCount
	}
	text.SetRank(rank)
	err = service.textRepo.Store(text)
	if err != nil {
		return err
	}

	channel := "personal#" + hash
	err = service.centrifugoClient.Publish(
		channel,
		map[string]interface{}{
			"hash":       hash,
			"rank":       text.GetRank(),
			"similarity": text.GetSimilarity(),
		},
	)
	if err != nil {
		return err
	}
	return service.eventDispatcher.Dispatch(appevent.NewRankCalculatedEvent(text.GetHash(), rank))
}
