package tradier

import (
	"cdr.dev/slog"
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/ppaanngggg/option-bot/pkg/market"
	"github.com/ppaanngggg/option-bot/pkg/utils"
	"golang.org/x/xerrors"
	"sort"
	"time"
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

// GetOptionChains refer to https://documentation.tradier.com/brokerage-api/markets/get-options-chains
func (t *Tradier) GetOptionChains(
	ctx context.Context, symbol string, expiration string,
) ([]market.Chain, error) {
	resp, err := t.client.R().
		SetContext(ctx).
		SetAuthToken(t.apiKey).
		SetHeader("Accept", "application/json").
		SetQueryParams(
			map[string]string{
				"symbol":     symbol,
				"expiration": expiration,
				"greeks":     "true",
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
	// Parse response body
	body := &struct {
		Options struct {
			Option []struct {
				Symbol         string  `json:"symbol"`
				OptionType     string  `json:"option_type"`     // call or put
				ExpirationDate string  `json:"expiration_date"` // YYYY-MM-DD
				Underlying     string  `json:"underlying"`      // SPX
				RootSymbol     string  `json:"root_symbol"`     // SPXW or SPX
				ContractSize   int     `json:"contract_size"`   // 100
				Strike         float64 `json:"strike"`
				Bid            float64 `json:"bid"`
				BidSize        int     `json:"bidsize"`
				BidDate        int64   `json:"bid_date"`
				Ask            float64 `json:"ask"`
				AskSize        int     `json:"asksize"`
				AskDate        int64   `json:"ask_date"`
				Greeks         struct {
					Delta     float64 `json:"delta"`
					Gamma     float64 `json:"gamma"`
					Theta     float64 `json:"theta"`
					Vega      float64 `json:"vega"`
					BidIV     float64 `json:"bid_iv"`
					MidIV     float64 `json:"mid_iv"`
					AskIV     float64 `json:"ask_iv"`
					UpdatedAt string  `json:"updated_at"` // YYYY-MM-DD HH:MM:SS
				} `json:"greeks"`
			} `json:"option"`
		} `json:"options"`
	}{}
	if err = json.Unmarshal(resp.Body(), body); err != nil {
		return nil, xerrors.New(err.Error())
	}
	// the datetime format is in EST
	timeLocation, err := time.LoadLocation("EST")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	// Group options by root_symbol
	chains := make(map[string]market.Chain)
	for _, opt := range body.Options.Option {
		chain, ok := chains[opt.RootSymbol]
		if !ok {
			chain = market.Chain{
				Symbol:         opt.RootSymbol,
				ExpirationDate: opt.ExpirationDate,
			}
		}
		// Append option to chain
		option := market.Option{
			StrikePrice: opt.Strike,
			BidPrice:    opt.Bid,
			BidSize:     opt.BidSize,
			BidAt:       opt.BidDate,
			AskPrice:    opt.Ask,
			AskSize:     opt.AskSize,
			AskAt:       opt.AskDate,
			QuoteAt:     time.Now().UnixMilli(),
			IV:          opt.Greeks.BidIV,
			Delta:       opt.Greeks.Delta,
			Gamma:       opt.Greeks.Gamma,
			Vega:        opt.Greeks.Vega,
			Theta:       opt.Greeks.Theta,
		}
		// parse greeks updated at to unix milli
		greeksUpdateAt, err := time.ParseInLocation(
			"2006-01-02 15:04:05", opt.Greeks.UpdatedAt, timeLocation,
		)
		if err != nil {
			return nil, xerrors.New(err.Error())
		}
		option.GreeksUpdatedAt = greeksUpdateAt.UnixMilli()
		if opt.OptionType == "call" {
			chain.Calls = append(chain.Calls, option)
		} else {
			chain.Puts = append(chain.Puts, option)
		}
		chains[opt.RootSymbol] = chain
	}
	// Sort calls and puts by strike price, and return
	rets := make([]market.Chain, 0, len(chains))
	for _, chain := range chains {
		chain.SortByStrikePrice()
		rets = append(rets, chain)
	}
	return rets, nil
}

// GetOptionExpirations refer to https://documentation.tradier.com/brokerage-api/markets/get-options-expirations
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
				"includeAllRoots": "true", // to deal with VIX/VIXW, SPX/SPXW, NDX/NDXP, RUT/RUTW
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
	// Parse response body
	body := &struct {
		Expirations struct {
			Date []string `json:"date"`
		} `json:"expirations"`
	}{}
	if err = json.Unmarshal(resp.Body(), body); err != nil {
		return nil, xerrors.New(err.Error())
	}
	// sort expirations by date
	sort.Strings(body.Expirations.Date)
	return body.Expirations.Date, nil
}
