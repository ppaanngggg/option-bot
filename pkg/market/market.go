package market

import (
	"context"
	"sort"

	marketv1 "github.com/ppaanngggg/option-bot/proto/gen/market/v1"
)

type Market interface {
	// Search returns a list of symbols that match the given query
	Search(ctx context.Context, query string) ([]*marketv1.Symbol, error)
	// GetOptionExpirations returns a list of expiration dates for the given symbol, e.g. "SPY"
	// return as an asc list of strings in the format "YYYY-MM-DD", e.g. ["2021-01-15", "2021-02-19", ...]
	GetOptionExpirations(ctx context.Context, symbol string) ([]string, error)
	// GetOptionChains returns a list of option chains for the given symbol and expiration date
	// Because some symbols have multiple option chains, for example, SPX may return one for SPXW and one for SPY
	GetOptionChains(
		ctx context.Context, symbol string, expiration string,
	) ([]*marketv1.Chain, error)
	// GetTodayTradePeriod returns the trading period for today
	GetTodayTradePeriod(ctx context.Context) (*marketv1.TradePeriod, error)
}

func SortByStrikePrice(c *marketv1.Chain) {
	// sort calls and puts by strike price
	sort.Slice(
		c.Calls, func(i, j int) bool {
			return c.Calls[i].Strike < c.Calls[j].Strike
		},
	)
	sort.Slice(
		c.Puts, func(i, j int) bool {
			return c.Puts[i].Strike < c.Puts[j].Strike
		},
	)
}
