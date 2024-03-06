package tradier

import (
	"cdr.dev/slog"
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/ppaanngggg/option-bot/pkg/market"
	"github.com/ppaanngggg/option-bot/pkg/utils"
	"golang.org/x/xerrors"
)

type Tradier struct {
	isLive bool
	apiKey string
	client *resty.Client
	logger slog.Logger
}

func NewTradier(isLive bool, apiKey string) *Tradier {
	tradier := &Tradier{
		isLive: isLive,
		apiKey: apiKey,
		client: resty.New(),
		logger: utils.DefaultLogger.With(slog.F("broker", "tradier")),
	}
	if tradier.isLive {
		tradier.client.SetBaseURL("https://api.tradier.com/v1/")
	} else {
		tradier.client.SetBaseURL("https://sandbox.tradier.com/v1/")
	}
	return tradier
}

func (t *Tradier) GetOptionChains(
	ctx context.Context, symbol string, expiration string,
) (*market.Option, error) {
	resp, err := t.client.R().
		SetContext(ctx).
		SetAuthToken(t.apiKey).
		SetHeader("Accept", "application/json").
		SetQueryParams(
			map[string]string{
				"symbol":     symbol,
				"expiration": expiration,
			},
		).Get("/markets/options/chains")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if resp.IsError() {
		return nil, xerrors.Errorf(
			"failed to get option chains, status: %s, body: %s",
			resp.Status(), resp.String(),
		)
	}
	t.logger.Info(ctx, resp.String())
	return nil, nil
}

func (t *Tradier) GetOptionExpirations(
	ctx context.Context, symbol string,
) ([]string, error) {
	resp, err := t.client.R().
		SetContext(ctx).
		SetAuthToken(t.apiKey).
		SetHeader("Accept", "application/json").
		SetQueryParams(
			map[string]string{
				"symbol":          symbol,
				"includeAllRoots": "true",
				"expirationType":  "true",
			},
		).
		Get("/markets/options/expirations")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if resp.IsError() {
		return nil, xerrors.Errorf(
			"failed to get option expirations, status: %s, body: %s",
			resp.Status(), resp.String(),
		)
	}
	// TODO: parse
	t.logger.Info(ctx, resp.String())
	return nil, nil
}
