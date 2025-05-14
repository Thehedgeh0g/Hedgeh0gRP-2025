package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CentrifugoClient interface {
	Publish(channel string, data interface{}) error
}

type centrifugoClient struct {
}

func NewCentrifugoClient() CentrifugoClient {
	return &centrifugoClient{}
}

func (cc *centrifugoClient) Publish(channel string, data interface{}) error {

	url := "http://centrifugo:8000/api/publish"

	payload := map[string]interface{}{
		"method":  "publish",
		"channel": channel,
		"data":    data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "apikey _salt")
	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
