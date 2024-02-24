package ibkr

import (
	"context"
	"testing"
)

func newIBKR() *IBKR {
	return NewIBKR("http://localhost:8000")
}

func TestIBKR_Login(t *testing.T) {
	ibkr := newIBKR()
	err := ibkr.Login(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
