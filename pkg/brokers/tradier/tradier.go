package tradier

import (
	"context"
	"sort"
	"time"

	"cdr.dev/slog"
	"github.com/go-resty/resty/v2"
	"github.com/ppaanngggg/option-bot/pkg/market"
	"github.com/ppaanngggg/option-bot/pkg/utils"
	marketv1 "github.com/ppaanngggg/option-bot/proto/gen/market/v1"
	"golang.org/x/xerrors"
)

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

var _ market.Market = (*Tradier)(nil)

type Tradier struct {
	isLive bool
	apiKey string
	client *resty.Client
	logger slog.Logger
}

// Search refer to https://documentation.tradier.com/brokerage-api/markets/get-lookup
func (t *Tradier) Search(ctx context.Context, query string) (
	[]*marketv1.Symbol, error,
) {
	body := &struct {
		Securities struct {
			Security []struct {
				Symbol      string `json:"symbol"`
				Exchange    string `json:"exchange"`
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"security"`
		} `json:"securities"`
	}{}
	resp, err := t.client.R().
		SetContext(ctx).
		SetAuthToken(t.apiKey).
		SetHeader("Accept", "application/json").
		SetQueryParams(
			map[string]string{
				"q": query,
			},
		).
		SetResult(body).
		Get("/markets/lookup")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if resp.IsError() {
		return nil, xerrors.Errorf(
			"failed to search symbols, status: %s, body: %s",
			resp.Status(), resp.String(),
		)
	}
	// Convert to marketv1.Symbol
	symbols := make([]*marketv1.Symbol, 0, len(body.Securities.Security))
	for _, sec := range body.Securities.Security {
		symbol := &marketv1.Symbol{
			Symbol:      sec.Symbol,
			Description: sec.Description,
		}
		switch sec.Type {
		case "stock":
			symbol.Type = marketv1.SymbolType_SYMBOL_TYPE_STOCK
		case "option":
			symbol.Type = marketv1.SymbolType_SYMBOL_TYPE_OPTION
		case "index":
			symbol.Type = marketv1.SymbolType_SYMBOL_TYPE_INDEX
		case "etf":
			symbol.Type = marketv1.SymbolType_SYMBOL_TYPE_ETF
		}
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

// GetTodayTradePeriod refer to https://documentation.tradier.com/brokerage-api/markets/get-calendar
func (t *Tradier) GetTodayTradePeriod(ctx context.Context) (
	*marketv1.TradePeriod, error,
) {
	body := &struct {
		Calendar struct {
			Month int `json:"month"`
			Year  int `json:"year"`
			Days  struct {
				Day []struct {
					Date      string `json:"date"`
					Status    string `json:"status"` // open or closed
					Premarket struct {
						Start string `json:"start"`
						End   string `json:"end"`
					} `json:"premarket"`
					Open struct {
						Start string `json:"start"` // 09:30
						End   string `json:"end"`   // 16:00
					} `json:"open"`
					Postmarket struct {
						Start string `json:"start"`
						End   string `json:"end"`
					} `json:"postmarket"`
				} `json:"day"`
			} `json:"days"`
		} `json:"calendar"`
	}{}
	resp, err := t.client.R().
		SetContext(ctx).
		SetAuthToken(t.apiKey).
		SetHeader("Accept", "application/json").
		SetResult(body).
		Get("/markets/calendar")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if resp.IsError() {
		return nil, xerrors.Errorf(
			"failed to get today trade period, status: %s, body: %s",
			resp.Status(), resp.String(),
		)
	}
	// get today's date in the format "YYYY-MM-DD" of New York timezone
	today := time.Now().In(utils.TZNewYork).Format("2006-01-02")
	// find today's trade period
	for _, day := range body.Calendar.Days.Day {
		if day.Date == today {
			if day.Status == "open" {
				startTime, err := time.ParseInLocation(
					"2006-04-02 15:04:05", today+" "+day.Open.Start+":00",
					utils.TZNewYork,
				)
				if err != nil {
					return nil, xerrors.New(err.Error())
				}
				endTime, err := time.ParseInLocation(
					"2006-04-02 15:04:05", today+" "+day.Open.End+":00",
					utils.TZNewYork,
				)
				if err != nil {
					return nil, xerrors.New(err.Error())
				}
				return &marketv1.TradePeriod{
					Date:    day.Date,
					IsOpen:  true,
					OpenAt:  startTime.UnixMilli(),
					CloseAt: endTime.UnixMilli(),
				}, nil
			} else {
				return &marketv1.TradePeriod{
					Date: day.Date,
				}, nil
			}
		}
	}
	return nil, xerrors.Errorf("today's trade period not found, today: %s", today)
}

// GetOptionChains refer to https://documentation.tradier.com/brokerage-api/markets/get-options-chains
func (t *Tradier) GetOptionChains(
	ctx context.Context, symbol string, expiration string,
) ([]*marketv1.Chain, error) {
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
		).
		SetResult(body).
		Get("/markets/options/chains")
	if err != nil {
		return nil, xerrors.New(err.Error())
	}
	if resp.IsError() {
		return nil, xerrors.Errorf(
			"failed to get option chains, status: %s, body: %s",
			resp.Status(), resp.String(),
		)
	}
	// Group options to chain by root_symbol
	chains := make(map[string]*marketv1.Chain)
	for _, opt := range body.Options.Option {
		// e.g. SPXW or SPX for underlying SPX
		chain, ok := chains[opt.RootSymbol]
		if !ok {
			chain = &marketv1.Chain{
				RootSymbol: opt.RootSymbol,
				Underlying: opt.Underlying,
				Expiration: opt.ExpirationDate,
			}
		}
		// Append option to chain
		option := &marketv1.Option{
			Symbol:  opt.Symbol,
			Strike:  opt.Strike,
			Bid:     opt.Bid,
			BidSize: int32(opt.BidSize),
			BidAt:   opt.BidDate,
			Ask:     opt.Ask,
			AskSize: int32(opt.AskSize),
			AskAt:   opt.AskDate,
			QuoteAt: time.Now().UnixMilli(),
			Iv:      opt.Greeks.BidIV,
			Delta:   opt.Greeks.Delta,
			Gamma:   opt.Greeks.Gamma,
			Vega:    opt.Greeks.Vega,
			Theta:   opt.Greeks.Theta,
		}
		// parse greeks updated at to unix milli
		greeksUpdateAt, err := time.ParseInLocation(
			"2006-01-02 15:04:05", opt.Greeks.UpdatedAt, utils.TZNewYork,
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
	rets := make([]*marketv1.Chain, 0, len(chains))
	for _, chain := range chains {
		market.SortByStrikePrice(chain)
		rets = append(rets, chain)
	}
	return rets, nil
}

// GetOptionExpirations refer to https://documentation.tradier.com/brokerage-api/markets/get-options-expirations
func (t *Tradier) GetOptionExpirations(
	ctx context.Context, symbol string,
) ([]string, error) {
	body := &struct {
		Expirations struct {
			Date []string `json:"date"`
		} `json:"expirations"`
	}{}
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
		SetResult(body).
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
	// sort expirations by date
	sort.Strings(body.Expirations.Date)
	return body.Expirations.Date, nil
}
