package ibkr

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/ppaanngggg/option-bot/pkg/utils"
)

type IBKR struct {
	host   string
	client *resty.Client
}

func NewIBKR(host string) *IBKR {
	ibkr := &IBKR{host: host, client: resty.New()}
	ibkr.client.SetBaseURL(ibkr.host)
	return ibkr
}

func (b *IBKR) Login(ctx context.Context) error {
	resp, err := b.client.R().SetContext(ctx).Post("/v1/api/login")
	if err != nil {
		return err
	}
	utils.DefaultLogger.Info(ctx, "ibkr login resp", "status", resp.Status())
	return nil
}
