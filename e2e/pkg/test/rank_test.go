package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"

	"e2etests/pkg/utils"
	"e2etests/pkg/utils/page"
)

func TestE2E(t *testing.T) {
	utils.RunTestForBrowser(t, "chrome", testCalculateTextRank)
	utils.RunTestForBrowser(t, "firefox", testCalculateTextRank)
}

func testCalculateTextRank(t *testing.T, driver selenium.WebDriver) {
	indexPage := page.Index{}
	indexPage.Init(driver)

	err := indexPage.OpenPage("/")
	assert.NoError(t, err, "Не удалось открыть главную страницу")

	expectedText := "A Б"
	expectedRank := 2.0 / 3.0
	expectedSimilarity := 1

	err = indexPage.InputText(expectedText)
	assert.NoError(t, err, "Не удалось ввести текст")

	err = indexPage.SelectRegion("RU")
	assert.NoError(t, err, "Не удалось выбрать регион")

	err = indexPage.SubmitText()
	assert.NoError(t, err, "Не удалось отправить текст")

	summaryPage := page.Summary{}
	summaryPage.Init(driver)

	time.Sleep(20 * time.Millisecond)
	actualRank, err := summaryPage.GetResultRank()
	assert.NoError(t, err)
	assert.Equal(t, expectedRank, actualRank, "Rank не совпадает")

	actualSimilarity, err := summaryPage.GetResultSimilarity()
	assert.NoError(t, err)
	assert.Equal(t, expectedSimilarity, actualSimilarity, "Similarity не совпадает")
}
