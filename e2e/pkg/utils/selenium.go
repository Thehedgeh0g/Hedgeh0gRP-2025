package utils

import (
	"e2etests/config"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

type SeleniumAdapter struct {
	driver  selenium.WebDriver
	baseURL string
}

func (s *SeleniumAdapter) Init(driver selenium.WebDriver) {
	s.driver = driver
	s.baseURL = config.BaseUrl
}

func (s *SeleniumAdapter) OpenPage(url string) error {
	return s.driver.Get(s.baseURL + url)
}

func (s *SeleniumAdapter) FindElement(by, value string) (selenium.WebElement, error) {
	return s.waitForElement(by, value, 10*time.Second)
}

func (s *SeleniumAdapter) waitForElement(by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var elem selenium.WebElement
	var err error

	startTime := time.Now()
	for time.Since(startTime) < timeout {
		elem, err = s.driver.FindElement(by, value)
		if err == nil {
			return elem, nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("element not found after %v: %v", timeout, err)
}
