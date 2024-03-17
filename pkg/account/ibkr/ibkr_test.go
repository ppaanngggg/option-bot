package ibkr

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newIBKR() *IBKR {
	return NewIBKR("http://localhost:8000")
}

func TestIBKR_Login(t *testing.T) {
	ibkr := newIBKR()
	err := ibkr.Login(context.Background())
	assert.NoError(t, err)
}
