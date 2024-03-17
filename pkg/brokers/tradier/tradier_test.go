package tradier

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTradier() *Tradier {
	return NewTradier(false, os.Getenv("UNITTEST_TRADIER_API_KEY"))
}

func TestTradier_Market(t *testing.T) {
	tradier := newTradier()
	ctx := context.Background()

	tradePeriod, err := tradier.GetTodayTradePeriod(ctx)
	assert.NoError(t, err)
	t.Logf("today: %v", tradePeriod.Date)
	if tradePeriod.IsOpen {
		t.Logf("trading period: %v - %v", tradePeriod.OpenAt, tradePeriod.CloseAt)
	} else {
		t.Logf("market is closed")
	}

	symbols, err := tradier.Search(ctx, "SPX")
	assert.NoError(t, err)
	assert.NotEmpty(t, symbols)
	t.Logf("number of symbols: %d", len(symbols))
	t.Logf("first symbol: %s", symbols[0].Symbol)

	exps, err := tradier.GetOptionExpirations(ctx, "SPX")
	assert.NoError(t, err)
	assert.NotEmpty(t, exps)
	t.Logf("number of expirations: %d", len(exps))
	t.Logf("first expiration: %s", exps[0])

	chains, err := tradier.GetOptionChains(ctx, "SPX", exps[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, chains)
	assert.Equal(t, "SPX", chains[0].Underlying)
	assert.Equal(t, exps[0], chains[0].Expiration)
	t.Logf("number of chains: %d", len(chains))
	t.Logf(
		"first chain root symbol: %s, underlying: %s, expiration: %s",
		chains[0].RootSymbol, chains[0].Underlying, chains[0].Expiration,
	)
	t.Logf(
		"first call option symbol: %s, put option symbol: %s",
		chains[0].Calls[0].Symbol, chains[0].Puts[0].Symbol,
	)
}
