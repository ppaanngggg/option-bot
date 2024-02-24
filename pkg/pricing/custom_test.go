package pricing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ppaanngggg/option-bot/pkg/market"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readOptionsDX() (*market.Chain, error) {
	file, err := os.Open("./test_data.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	_, err = reader.Read() // skip header
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	chain := &market.Chain{
		Symbol: "SPY",
	}
	for _, record := range records {
		callSize := record[16]
		callSizeSplit := strings.Split(callSize, "x")
		callBidSize, err := strconv.Atoi(strings.TrimSpace(callSizeSplit[0]))
		if err != nil {
			return nil, err
		}
		callAskSize, err := strconv.Atoi(strings.TrimSpace(callSizeSplit[1]))
		if err != nil {
			return nil, err
		}
		callBid := record[17]
		callBidPrice, err := strconv.ParseFloat(callBid, 64)
		if err != nil {
			return nil, err
		}
		callAsk := record[18]
		callAskPrice, err := strconv.ParseFloat(callAsk, 64)
		if err != nil {
			return nil, err
		}
		strike := record[19]
		strikePrice, err := strconv.ParseFloat(strike, 64)
		if err != nil {
			return nil, err
		}
		putBid := record[20]
		putBidPrice, err := strconv.ParseFloat(putBid, 64)
		if err != nil {
			return nil, err
		}
		putAsk := record[21]
		putAskPrice, err := strconv.ParseFloat(putAsk, 64)
		if err != nil {
			return nil, err
		}
		putSize := record[22]
		putSizeSplit := strings.Split(putSize, "x")
		putBidSize, err := strconv.Atoi(strings.TrimSpace(putSizeSplit[0]))
		if err != nil {
			return nil, err
		}
		putAskSize, err := strconv.Atoi(strings.TrimSpace(putSizeSplit[1]))
		if err != nil {
			return nil, err
		}
		chain.Calls = append(
			chain.Calls, market.Option{
				StrikePrice: strikePrice,
				BidPrice:    callBidPrice,
				BidSize:     callBidSize,
				AskPrice:    callAskPrice,
				AskSize:     callAskSize,
			},
		)
		chain.Puts = append(
			chain.Puts, market.Option{
				StrikePrice: strikePrice,
				BidPrice:    putBidPrice,
				BidSize:     putBidSize,
				AskPrice:    putAskPrice,
				AskSize:     putAskSize,
			},
		)
	}
	return chain, nil
}

func TestCustomOptionsDX(t *testing.T) {
	chain, err := readOptionsDX()
	if err != nil {
		t.Fatal(err)
	}
	if len(chain.Calls) != len(chain.Puts) {
		t.Fatalf(
			"calls and puts length mismatch: %d != %d", len(chain.Calls),
			len(chain.Puts),
		)
	}
	chain.SortByStrikePrice()
	priceDistributions, err := PredictPriceDistribution(chain)
	if err != nil {
		t.Fatal(err)
	}
	for _, pd := range priceDistributions {
		fmt.Printf("price: %0.4f, prob: %0.4f\n", pd.Price, pd.Prob)
	}
	for i := range chain.Calls {
		call := chain.Calls[i]
		callDelta := 0.0
		for _, pd := range priceDistributions {
			if pd.Price > call.StrikePrice {
				callDelta += pd.Prob
			}
		}
		put := chain.Puts[i]
		putDelta := 0.0
		for _, pd := range priceDistributions {
			if pd.Price < put.StrikePrice {
				putDelta += pd.Prob
			}
		}
		fmt.Printf(
			"strike: %0.4f, callDelta: %0.4f, putDelta: %0.4f\n",
			call.StrikePrice, callDelta, putDelta,
		)
	}
}

type tradierOptions struct {
	Options struct {
		Option []struct {
			Strike     float64 `json:"strike"`
			Bid        float64 `json:"bid"`
			BidSize    int     `json:"bidsize"`
			Ask        float64 `json:"ask"`
			AskSize    int     `json:"asksize"`
			OptionType string  `json:"option_type"`
		} `json:"option"`
	} `json:"options"`
}

func readTradier() (*market.Chain, error) {
	chain := &market.Chain{
		Symbol: "SPY",
	}
	file, err := os.Open("./test_data.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var tradier tradierOptions
	err = decoder.Decode(&tradier)
	if err != nil {
		return nil, err
	}
	for _, option := range tradier.Options.Option {
		switch option.OptionType {
		case "call":
			chain.Calls = append(
				chain.Calls, market.Option{
					StrikePrice: option.Strike,
					BidPrice:    option.Bid,
					BidSize:     option.BidSize,
					AskPrice:    option.Ask,
					AskSize:     option.AskSize,
				},
			)
		case "put":
			chain.Puts = append(
				chain.Puts, market.Option{
					StrikePrice: option.Strike,
					BidPrice:    option.Bid,
					BidSize:     option.BidSize,
					AskPrice:    option.Ask,
					AskSize:     option.AskSize,
				},
			)
		}
	}
	return chain, nil
}

func TestCustomTradier(t *testing.T) {
	chain, err := readTradier()
	if err != nil {
		t.Fatal(err)
	}
	if len(chain.Calls) != len(chain.Puts) {
		t.Fatalf(
			"calls and puts length mismatch: %d != %d", len(chain.Calls),
			len(chain.Puts),
		)
	}
	chain.SortByStrikePrice()
	priceDistributions, err := PredictPriceDistribution(chain)
	if err != nil {
		t.Fatal(err)
	}
	for _, pd := range priceDistributions {
		fmt.Printf("price: %0.4f, prob: %0.4f\n", pd.Price, pd.Prob)
	}
	for i := range chain.Calls {
		call := chain.Calls[i]
		callDelta := 0.0
		for _, pd := range priceDistributions {
			if pd.Price > call.StrikePrice {
				callDelta += pd.Prob
			}
		}
		put := chain.Puts[i]
		putDelta := 0.0
		for _, pd := range priceDistributions {
			if pd.Price < put.StrikePrice {
				putDelta += pd.Prob
			}
		}
		fmt.Printf(
			"strike: %0.4f, callDelta: %0.4f, putDelta: %0.4f\n",
			call.StrikePrice, callDelta, putDelta,
		)
	}
}
