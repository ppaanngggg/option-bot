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
	Type  string  `json:"type"`  // exact, nearest, at_least, at_most
	Delta float64 `json:"delta"` // choose by delta

	Price        float64 `json:"price"`         // choose by price
	StrikeOffset float64 `json:"strike_offset"` // choose by strike
}

type DTE struct {
	Type string `json:"type"` // exact, nearest, at_least, at_most
	Days int    `json:"day"`  // choose by days

}

type Leg struct {
	Underlying string `json:"underlying"`
	Action     string `json:"action"`   // buy/sell
	Type       string `json:"type"`     // call/put
	Quantity   int    `json:"quantity"` // the final leg's size is position's size * quantity
	Strike     `json:"strike"`
}

type Entry struct {
}

type Exit struct {
}

type Setting struct {
	Legs  []Leg `json:"legs"`
	Entry `json:"entry"`
	Exit  `json:"exit"`
}

type Position struct {
}

type Bot struct {
	Setting   `json:"setting"`
	Name      string `json:"name"`
	AutoOpen  bool   `json:"auto_open"`
	AutoClose bool   `json:"auto_close"`

	Positions []Position `json:"positions"`
}
