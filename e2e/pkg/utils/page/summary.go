package page

import (
	"e2etests/pkg/utils"
	"strconv"

	"github.com/tebeka/selenium"
)

const (
	resultRankCSSSelector       = ".container p:nth-of-type(1) .result-value"
	resultSimilarityCSSSelector = ".container p:nth-of-type(2) .result-value"
)

type Summary struct {
	utils.SeleniumAdapter
}

func (s *Summary) GetResultRank() (float64, error) {
	val, err := s.getResultValue(resultRankCSSSelector)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

func (s *Summary) GetResultSimilarity() (int, error) {
	val, err := s.getResultValue(resultSimilarityCSSSelector)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func (s *Summary) getResultValue(selector string) (string, error) {
	elem, err := s.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		return "", err
	}

	text, err := elem.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}
