package pricing

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readOptionsDX() (*Chain, error) {
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
	chain := Chain{
		Symbol:         "SPY",
		QuoteUnixTime:  1701464400,
		ExpireUnixTime: 1701723600,
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
			chain.Calls, Option{
				StrikePrice: strikePrice,
				BidPrice:    callBidPrice,
				BidSize:     callBidSize,
				AskPrice:    callAskPrice,
				AskSize:     callAskSize,
			},
		)
		chain.Puts = append(
			chain.Puts, Option{
				StrikePrice: strikePrice,
				BidPrice:    putBidPrice,
				BidSize:     putBidSize,
				AskPrice:    putAskPrice,
				AskSize:     putAskSize,
			},
		)
	}
	return &chain, nil
}

func TestCustom(t *testing.T) {
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
	priceDistributions, err := chain.PredictPriceDistributionByCalls()
	if err != nil {
		t.Fatal(err)
	}
	for _, pd := range priceDistributions {
		fmt.Printf("price: %0.4f, prob: %0.4f\n", pd.Price, pd.Prob)
	}
	for _, call := range chain.Calls {
		delta := 0.0
		for _, pd := range priceDistributions {
			if pd.Price > call.StrikePrice {
				delta += pd.Prob
			}
		}
		fmt.Printf("strike: %0.4f, delta: %0.4f\n", call.StrikePrice, delta)
	}
}
