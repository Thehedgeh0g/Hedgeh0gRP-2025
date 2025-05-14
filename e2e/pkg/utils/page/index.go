package page

import (
	"e2etests/pkg/utils"
	"fmt"

	"github.com/tebeka/selenium"
)

const (
	textAreaCSSSelector     = "textarea[name='text']"
	selectRegionCSSSelector = "select[name='region']"
	optionRegionXPath       = ".//option[@value='%s']"
	submitTextCSSSelector   = "button[type='submit']"
)

type Index struct {
	utils.SeleniumAdapter
}

func (i *Index) InputText(text string) error {
	textArea, err := i.FindElement(selenium.ByCSSSelector, textAreaCSSSelector)
	if err != nil {
		return err
	}

	if err := textArea.Clear(); err != nil {
		return err
	}

	return textArea.SendKeys(text)
}

func (i *Index) SelectRegion(region string) error {
	selectElement, err := i.FindElement(selenium.ByCSSSelector, selectRegionCSSSelector)
	if err != nil {
		return err
	}

	option, err := selectElement.FindElement(selenium.ByXPATH, fmt.Sprintf(optionRegionXPath, region))
	if err != nil {
		return err
	}

	return option.Click()
}

func (i *Index) SubmitText() error {
	submitElement, err := i.FindElement(selenium.ByCSSSelector, submitTextCSSSelector)
	if err != nil {
		return err
	}

	return submitElement.Click()
}
