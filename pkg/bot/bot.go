package bot

const (
	BUY  = "buy"
	SELL = "sell"
	CALL = "call"
	PUT  = "put"

	EXACT    = "exact"
	NEAREST  = "nearest"
	AT_LEAST = "at_least"
	AT_MOST  = "at_most"
)

type Strike struct {
	Type       string  `json:"type"`  // exact, nearest, at_least, at_most
	Delta      float64 `json:"delta"` // choose by delta
	DeltaRange struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"delta_range"` // filter by delta range
	Price      float64 `json:"price"` // choose by price
	PriceRange struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"price_range"` // filter by price range
	StrikeOffset float64 `json:"strike_offset"` // choose by strike
}

type DTE struct {
	Type      string `json:"type"` // exact, nearest, at_least, at_most
	Days      int    `json:"day"`  // choose by days
	DaysRange struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"days_range"` // filter by days range
}

type Leg struct {
	Underlying string `json:"underlying"`
	Action     string `json:"action"`   // buy/sell
	Type       string `json:"type"`     // call/put
	Quantity   int    `json:"quantity"` // the final leg's size is position's size * quantity
	Strike     `json:"strike"`
	DTE        `json:"dte"`
}

type Entry struct {
	Weekdays struct {
		Monday    bool `json:"monday"`
		Tuesday   bool `json:"tuesday"`
		Wednesday bool `json:"wednesday"`
		Thursday  bool `json:"thursday"`
		Friday    bool `json:"friday"`
	} `json:"weekdays"`
	Time struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
		Second int `json:"second"`
	} `json:"time"`
}

type Exit struct {
	StopWin  float64 `json:"stop_win"`
	StopLoss float64 `json:"stop_loss"`
	Time     struct {
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
		Second int `json:"second"`
	} `json:"time"`
}

type Allocation struct {
	// constant size, such as 1, 2, 3, 4, 5
	Constant int `json:"constant"`
	// size to make sure the max risk is less than max_risk, such as $1000
	MaxRisk float64 `json:"max_risk"`
}

type Setting struct {
	Legs       []Leg      `json:"legs"`
	Allocation Allocation `json:"allocation"`
	Entry      Entry      `json:"entry"`
	Exit       Exit       `json:"exit"`
}

type Bot struct {
	Name            string  `json:"name"`
	Setting         Setting `json:"setting"`
	EnableAutoOpen  bool    `json:"enable_auto_open"`
	EnableAutoClose bool    `json:"enable_auto_close"`
}
