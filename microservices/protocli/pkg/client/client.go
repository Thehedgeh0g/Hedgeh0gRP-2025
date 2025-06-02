package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrKeyNotFound = errors.New("Key not found")

type ProtoKeyClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProtoKeyClient(baseURL string) *ProtoKeyClient {
	return &ProtoKeyClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *ProtoKeyClient) Set(key string, value string) error {
	data := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	body, _ := json.Marshal(data)
	fmt.Println(c.baseURL+"/set", "application/json")
	resp, err := c.httpClient.Post(c.baseURL+"/set", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set failed: %s", resp.Status)
	}
	return nil
}

func (c *ProtoKeyClient) Get(key string) (string, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/get?key=" + url.QueryEscape(key))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var result struct {
			Value string `json:"value"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}
		return result.Value, nil

	case http.StatusBadRequest:
		return "", errors.New("invalid key")

	case http.StatusNotFound:
		return "", ErrKeyNotFound

	default:
		return "", fmt.Errorf("get failed: %s", resp.Status)
	}
}

func (c *ProtoKeyClient) Keys(prefix string) ([]string, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/keys?prefix=" + url.QueryEscape(prefix))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("invalid prefix")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keys failed: %s", resp.Status)
	}

	var result struct {
		Keys []string `json:"keys"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Keys, nil
}
