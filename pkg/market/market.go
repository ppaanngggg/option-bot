package market

import (
	"context"
	"sort"
	"time"
)

type Market interface {
	GetOptionChains(ctx context.Context, symbol string, expiration string) (*Option, error)
	GetOptionExpirations(ctx context.Context, symbol string) ([]string, error)
}

type Chain struct {
	Symbol     string
	QuoteTime  time.Time
	ExpireTime time.Time
	IV         float64
	Calls      []Option
	Puts       []Option
}

type Option struct {
	StrikePrice float64
	BidPrice    float64
	BidSize     int
	AskPrice    float64
	AskSize     int

	IV    float64
	Delta float64
	Gamma float64
	Vega  float64
	Theta float64
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
