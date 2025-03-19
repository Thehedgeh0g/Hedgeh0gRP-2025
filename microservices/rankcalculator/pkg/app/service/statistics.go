package service

import (
	"regexp"

	"rankcalculator/pkg/app/model"
)

type StatisticsService interface {
	CalculateRank(hash string) error
}

func NewStatisticsService(textRepo model.TextRepository) StatisticsService {
	return &statisticsService{
		textRepo: textRepo,
	}
}

type statisticsService struct {
	textRepo model.TextRepository
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
	return service.textRepo.Store(text)
}
