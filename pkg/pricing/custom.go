package pricing

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
