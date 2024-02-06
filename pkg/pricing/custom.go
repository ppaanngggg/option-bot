package pricing

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
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

func (c *Chain) findMinStrikeDiff() float64 {
	minStrikeDiff := math.MaxFloat64
	for i := 1; i < len(c.Calls); i++ {
		strikeDiff := c.Calls[i].StrikePrice - c.Calls[i-1].StrikePrice
		if strikeDiff < minStrikeDiff {
			minStrikeDiff = strikeDiff
		}
	}
	for i := 1; i < len(c.Puts); i++ {
		strikeDiff := c.Puts[i].StrikePrice - c.Puts[i-1].StrikePrice
		if strikeDiff < minStrikeDiff {
			minStrikeDiff = strikeDiff
		}
	}
	return minStrikeDiff
}

func (c *Chain) initPriceDistribution(minStrikeDiff float64) ([]float64, []float64) {
	minStrike := math.Min(c.Calls[0].StrikePrice, c.Puts[0].StrikePrice)
	maxStrike := math.Max(
		c.Calls[len(c.Calls)-1].StrikePrice, c.Puts[len(c.Puts)-1].StrikePrice,
	)
	price := minStrike - minStrikeDiff/2.0
	prices := make([]float64, 0)
	probs := make([]float64, 0)
	prices = append(prices, price)
	probs = append(probs, 0.0)
	for price < maxStrike {
		price += minStrikeDiff
		prices = append(prices, price)
		probs = append(probs, 0.0)
	}
	return prices, probs
}

func getRealPrice(option Option) float64 {
	var price float64
	if option.BidSize == 0 && option.AskSize == 0 {
		return 0
	}
	if option.BidPrice > 0 && option.AskPrice > 0 {
		price = (option.BidPrice + option.AskPrice) / 2.0
	}
	if option.BidPrice > 0 && option.AskPrice == 0 {
		price = option.BidPrice
	}
	if option.BidPrice == 0 && option.AskPrice > 0 {
		price = option.AskPrice
	}
	return price
}

func solvePriceDistribution(
	prices []float64, probs []float64, A []float64, b []float64, w []float64,
) ([]PriceDistribution, error) {
	// create matrix A and vector b
	Amat := mat.NewDense(len(A)/len(prices), len(prices), A)
	bmat := mat.NewDense(len(b), 1, b)
	wmat := mat.NewDense(len(w), 1, w)
	// calculate loss func
	lossFunc := func(x []float64) float64 {
		xmat := mat.NewDense(len(x), 1, x)
		// pow 2
		xmat.MulElem(xmat, xmat)
		// multiply A and x
		var tmp mat.Dense
		tmp.Mul(Amat, xmat)
		// subtract b
		tmp.Sub(&tmp, bmat)
		// pow 2
		tmp.MulElem(&tmp, &tmp)
		// multiply w
		tmp.MulElem(&tmp, wmat)
		// sum
		loss := mat.Sum(&tmp)
		return loss
	}
	gradFunc := func(grad, x []float64) {
		// TODO: speed up
		fd.Gradient(grad, lossFunc, x, nil)
	}
	// minimize loss func
	p := optimize.Problem{
		Func: lossFunc,
		Grad: gradFunc,
	}
	result, err := optimize.Minimize(
		p, probs, &optimize.Settings{Converger: &MyConverge{}}, &optimize.LBFGS{},
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

func (c *Chain) PredictPriceDistribution() ([]PriceDistribution, error) {
	prices, probs := c.initPriceDistribution(c.findMinStrikeDiff())
	// build matrix A and vector b
	A := make([]float64, 0)
	b := make([]float64, 0)
	w := make([]float64, 0)
	for _, call := range c.Calls {
		realPrice := getRealPrice(call)
		if realPrice == 0 {
			continue
		}
		b = append(b, realPrice)
		w = append(w, float64(call.BidSize+call.AskSize))
		for _, price := range prices {
			if price > call.StrikePrice {
				A = append(A, price-call.StrikePrice)
			} else {
				A = append(A, 0.0)
			}
		}
	}
	for _, put := range c.Puts {
		realPrice := getRealPrice(put)
		if realPrice == 0 {
			continue
		}
		b = append(b, realPrice)
		w = append(w, float64(put.BidSize+put.AskSize))
		for _, price := range prices {
			if price < put.StrikePrice {
				A = append(A, put.StrikePrice-price)
			} else {
				A = append(A, 0.0)
			}
		}
	}
	return solvePriceDistribution(prices, probs, A, b, w)
}
