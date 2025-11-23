package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPProvider fetches price from external HTTP endpoint.
// It expects JSON with field "price" (float).
type HTTPProvider struct {
	Client  *http.Client
	BaseURL string
	Path    string
}

type priceResponse struct {
	Price float64 `json:"price"`
}

func NewHTTPProvider(baseURL, path string, timeout time.Duration) *HTTPProvider {
	if timeout == 0 {
		timeout = 3 * time.Second
	}
	return &HTTPProvider{
		Client:  &http.Client{Timeout: timeout},
		BaseURL: baseURL,
		Path:    path,
	}
}

func (h *HTTPProvider) GetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("%s%s?symbol=%s", h.BaseURL, h.Path, symbol)
	resp, err := h.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http provider status %d", resp.StatusCode)
	}
	var pr priceResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return 0, err
	}
	if pr.Price <= 0 {
		return 0, fmt.Errorf("invalid price")
	}
	return pr.Price, nil
}
