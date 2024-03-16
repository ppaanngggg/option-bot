package ibkr

import (
	"context"

	"cdr.dev/slog"
	"github.com/go-resty/resty/v2"
	"github.com/ppaanngggg/option-bot/pkg/utils"
	"golang.org/x/xerrors"
)

type IBKR struct {
	host   string
	client *resty.Client
	logger slog.Logger
}

func NewIBKR(host string) *IBKR {
	ibkr := &IBKR{
		host:   host,
		client: resty.New(),
		logger: utils.DefaultLogger.With(slog.F("broker", "ibkr")),
	}
	ibkr.client.SetBaseURL(ibkr.host)
	return ibkr
}

// Login Please refer to https://github.com/ppaanngggg/ib-cp-server
func (i *IBKR) Login(ctx context.Context) error {
	resp, err := i.client.R().SetContext(ctx).Post("/v1/api/login")
	if err != nil {
		return xerrors.New(err.Error())
	}
	if resp.IsError() {
		return xerrors.Errorf(
			"failed to login, status: %s, body: %s", resp.Status(), resp.String(),
		)
	}
	return nil
}

// TODO: implement the rest of the market interface
