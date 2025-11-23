package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TwelveDataProvider fetches price from TwelveData price endpoint.
// Endpoint: https://api.twelvedata.com/price?symbol=BTC/USD&apikey=XXX
type TwelveDataProvider struct {
	Client  *http.Client
	BaseURL string
	Path    string
	APIKey  string
}

type twelvePriceResp struct {
	Price   string `json:"price"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func NewTwelveDataProvider(baseURL, path, apiKey string, timeout time.Duration) *TwelveDataProvider {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	if baseURL == "" {
		baseURL = "https://api.twelvedata.com"
	}
	if path == "" {
		path = "/price"
	}
	return &TwelveDataProvider{
		Client:  &http.Client{Timeout: timeout},
		BaseURL: baseURL,
		Path:    path,
		APIKey:  apiKey,
	}
}

func (t *TwelveDataProvider) GetPrice(symbol string) (float64, error) {
	if t.APIKey == "" {
		return 0, fmt.Errorf("twelvedata apiKey required")
	}
	u, _ := url.Parse(t.BaseURL)
	u.Path = t.Path
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("apikey", t.APIKey)
	u.RawQuery = q.Encode()

	resp, err := t.Client.Get(u.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("twelvedata status %d", resp.StatusCode)
	}
	var pr twelvePriceResp
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return 0, err
	}
	if pr.Status == "error" || pr.Code != 0 || pr.Message != "" {
		return 0, fmt.Errorf("twelvedata error: %s", pr.Message)
	}
	if pr.Price == "" {
		return 0, fmt.Errorf("twelvedata empty price")
	}
	price, err := parseFloat(pr.Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
