package pricing

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
)

type Option struct {
	StrikePrice float64
	BidPrice    float64
	BidSize     int
	AskPrice    float64
	AskSize     int

	IV    float64
	Delta float64
	Gamma float64
	Vega  float64
	Theta float64
}

type Chain struct {
	Symbol         string
	QuoteUnixTime  int // in seconds
	ExpireUnixTime int // in seconds
	IV             float64
	Calls          []Option
	Puts           []Option
}

type PriceDistribution struct {
	Price float64
	Prob  float64
}

func (c *Chain) SortByStrikePrice() {
	// sort calls and puts by strike price
	sort.Slice(
		c.Calls, func(i, j int) bool {
			return c.Calls[i].StrikePrice < c.Calls[j].StrikePrice
		},
	)
	sort.Slice(
		c.Puts, func(i, j int) bool {
			return c.Puts[i].StrikePrice < c.Puts[j].StrikePrice
		},
	)
}

type MyConverge struct {
	lastX          []float64
	diffThreshold  float64
	count          int
	countThreshold int
}

func (m *MyConverge) Init(dim int) {
	if m.diffThreshold == 0 {
		m.diffThreshold = 1e-5
	}
	if m.countThreshold == 0 {
		m.countThreshold = 10
	}
}

func (m *MyConverge) Converged(loc *optimize.Location) optimize.Status {
	maxDiff := 0.0
	for i, x := range loc.X {
		diff := math.Abs(x - m.lastX[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	m.lastX = loc.X
	if maxDiff < m.diffThreshold {
		m.count++
	} else {
		m.count = 0
	}
	if m.count > m.countThreshold {
		return optimize.FunctionConvergence
	}
	return optimize.NotTerminated
}

func (c *Chain) PredictPriceDistributionByCalls() ([]PriceDistribution, error) {
	// find minimum strike diff of calls
	minStrikeDiff := math.MaxFloat64
	for i := 1; i < len(c.Calls); i++ {
		strikeDiff := c.Calls[i].StrikePrice - c.Calls[i-1].StrikePrice
		if strikeDiff < minStrikeDiff {
			minStrikeDiff = strikeDiff
		}
	}
	// create init price distribution
	price := c.Calls[0].StrikePrice - minStrikeDiff/2.0
	prices := make([]float64, 0)
	probs := make([]float64, 0)
	for range c.Calls {
		price += minStrikeDiff
		prices = append(prices, price)
		probs = append(probs, 0.0)
	}
	// calculate loss func
	lossFunc := func(x []float64) float64 {
		loss := 0.0
		for _, call := range c.Calls {
			var optionPrice, shouldPrice float64
			if call.BidSize == 0 && call.AskSize == 0 {
				continue
			}
			if call.BidPrice > 0 && call.AskPrice > 0 {
				optionPrice = (call.BidPrice + call.AskPrice) / 2.0
			}
			if call.BidPrice > 0 && call.AskPrice == 0 {
				optionPrice = call.BidPrice
			}
			if call.BidPrice == 0 && call.AskPrice > 0 {
				optionPrice = call.AskPrice
			}
			for i, price := range prices {
				prob := math.Pow(x[i], 2)
				if price > call.StrikePrice {
					shouldPrice += prob * (price - call.StrikePrice)
				}
			}
			loss += math.Pow(optionPrice-shouldPrice, 2)
		}
		return loss
	}
	// minimize loss func
	p := optimize.Problem{
		Func: lossFunc,
		Grad: func(grad, x []float64) {
			fd.Gradient(grad, lossFunc, x, nil)
		},
	}
	result, err := optimize.Minimize(
		p, probs, &optimize.Settings{Converger: &MyConverge{lastX: probs}}, nil,
	)
	if err != nil {
		return nil, err
	}
	if err = result.Status.Err(); err != nil {
		return nil, err
	}
	fmt.Printf("result.Status: %v\n", result.Status)
	fmt.Printf("result.X: %0.4g\n", result.X)
	fmt.Printf("result.F: %0.4g\n", result.F)
	fmt.Printf("result.Stats.FuncEvaluations: %d\n", result.Stats.FuncEvaluations)
	// create price distribution
	priceDistributions := make([]PriceDistribution, 0)
	for i, price := range prices {
		prob := result.X[i]
		priceDistributions = append(
			priceDistributions,
			PriceDistribution{Price: price, Prob: math.Pow(prob, 2)},
		)
	}
	return priceDistributions, nil
}
