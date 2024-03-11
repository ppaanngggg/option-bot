package market

import (
	"context"
	"sort"
)

type Market interface {
	// GetOptionExpirations returns a list of expiration dates for the given symbol, e.g. "SPY"
	// return as an asc list of strings in the format "YYYY-MM-DD", e.g. ["2021-01-15", "2021-02-19", ...]
	GetOptionExpirations(ctx context.Context, symbol string) ([]string, error)
	// GetOptionChains returns a list of option chains for the given symbol and expiration date
	// Because some symbols have multiple option chains, for example, SPX may return one for SPXW and one for SPY
	GetOptionChains(ctx context.Context, symbol string, expiration string) ([]Chain, error)
}

type Chain struct {
	Symbol         string   `json:"symbol"`
	Underlying     string   `json:"underlying"`
	ExpirationDate string   `json:"expiration_date"`
	Calls          []Option `json:"calls"`
	Puts           []Option `json:"puts"`
}

type Option struct {
	StrikePrice float64 `json:"strike_price"`
	BidPrice    float64 `json:"bid_price"`
	BidSize     int     `json:"bid_size"`
	BidAt       int64   `json:"bid_at"` // Unix timestamp in milliseconds
	AskPrice    float64 `json:"ask_price"`
	AskSize     int     `json:"ask_size"`
	AskAt       int64   `json:"ask_at"`   // Unix timestamp in milliseconds
	QuoteAt     int64   `json:"quote_at"` // Unix timestamp in milliseconds

	IV              float64 `json:"iv"`
	Delta           float64 `json:"delta"`
	Gamma           float64 `json:"gamma"`
	Vega            float64 `json:"vega"`
	Theta           float64 `json:"theta"`
	GreeksUpdatedAt int64   `json:"greeks_updated_at"` // Unix timestamp in milliseconds
}

func (c *Chain) SortByStrikePrice() {
	// sort calls and puts by strike price
	sort.Slice(
		c.Calls, func(i, j int) bool {
			return c.Calls[i].StrikePrice < c.Calls[j].StrikePrice
		},
	)
	sort.Slice(
		c.Puts, func(i, j int) bool {
			return c.Puts[i].StrikePrice < c.Puts[j].StrikePrice
		},
	)
}
