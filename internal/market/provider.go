package market

// PriceProvider abstracts price retrieval for a symbol.
type PriceProvider interface {
	GetPrice(symbol string) (float64, error)
}
