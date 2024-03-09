package tradier

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func newTradier() *Tradier {
	return NewTradier(false, os.Getenv("UNITTEST_TRADIER_API_KEY"))
}

func TestTradier_Market(t *testing.T) {
	tradier := newTradier()
	exps, err := tradier.GetOptionExpirations(context.Background(), "SPX")
	assert.NoError(t, err)
	assert.NotEmpty(t, exps)
	for _, exp := range exps {
		println("SPX", exp)
	}

	chains, err := tradier.GetOptionChains(context.Background(), "SPX", exps[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, chains)
}
