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

func TestTradier_GetOptionExpirations(t *testing.T) {
	tradier := newTradier()
	exps, err := tradier.GetOptionExpirations(context.Background(), "SPX")
	assert.NoError(t, err)
	println(exps)
}

func TestTradier_GetOptionChains(t *testing.T) {
	tradier := newTradier()
	chains, err := tradier.GetOptionChains(context.Background(), "SPX", "2024-03-15")
	assert.NoError(t, err)
	println(chains)
}
