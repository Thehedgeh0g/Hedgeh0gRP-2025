package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tebeka/selenium"
)

func RunTestForBrowser(t *testing.T, browserName string, testFunc func(*testing.T, selenium.WebDriver)) {
	t.Helper()
	t.Run(browserName, func(t *testing.T) {
		caps := selenium.Capabilities{
			"browserName": browserName,
		}
		driver, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
		if !assert.NoError(t, err, "Failed to start "+browserName+" session") {
			return
		}
		defer driver.Quit()
		testFunc(t, driver)
	})
}
