package cbr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CBRClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCBRClient() *CBRClient {
	return &CBRClient{
		baseURL: "https://www.cbr-xml-daily.ru",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CBRResponse struct {
	Date         time.Time                   `json:"Date"`         // Для будущих обновлений.
	PreviousDate time.Time                   `json:"PreviousDate"` // Для будущих обновлений.
	PreviousURL  string                      `json:"PreviousURL"`  // Для будущих обновлений.
	Timestamp    time.Time                   `json:"Timestamp"`
	Valute       map[string]*CurrencyRespose `json:"Valute"`
}

type CurrencyRespose struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  float64 `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

func (c *CBRClient) GetDailyRates(ctx context.Context) (*CBRResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/daily_json.js", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var cbrResponse CBRResponse
	if err := json.NewDecoder(res.Body).Decode(&cbrResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("The current course has been received")
	return &cbrResponse, nil
}
