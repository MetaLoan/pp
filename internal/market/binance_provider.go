package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// BinanceProvider fetches ticker price from Binance public API.
// Response: { "symbol": "BTCUSDT", "price": "123.45" }
type BinanceProvider struct {
	Client  *http.Client
	BaseURL string
	Path    string
}

type binanceResp struct {
	Price string `json:"price"`
}

func NewBinanceProvider(baseURL, path string, timeout time.Duration) *BinanceProvider {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	if baseURL == "" {
		baseURL = "https://api.binance.com"
	}
	if path == "" {
		path = "/api/v3/ticker/price"
	}
	return &BinanceProvider{
		Client:  &http.Client{Timeout: timeout},
		BaseURL: baseURL,
		Path:    path,
	}
}

func (b *BinanceProvider) GetPrice(symbol string) (float64, error) {
	binanceSymbol := strings.ReplaceAll(symbol, "/", "")
	url := fmt.Sprintf("%s%s?symbol=%s", b.BaseURL, b.Path, binanceSymbol)
	resp, err := b.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("binance status %d", resp.StatusCode)
	}
	var br binanceResp
	if err := json.NewDecoder(resp.Body).Decode(&br); err != nil {
		return 0, err
	}
	price, err := parseFloat(br.Price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
