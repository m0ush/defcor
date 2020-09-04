package iex

// Stock represents a single Stock
type Stock struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	IexID    string `json:"iexId"`
	Region   string `json:"region"`
	Currency string `json:"currency"`
	IsActive bool   `json:"isEnabled"`
	Figi     string `json:"figi"`
	Cik      string `json:"cik"`
}

// Prices type models daily price data
type Prices struct {
	Date    string  `json:"date"`
	Uopen   float64 `json:"uOpen"`
	Uhigh   float64 `json:"uHigh"`
	Ulow    float64 `json:"uLow"`
	Uclose  float64 `json:"uClose"`
	Uvolume int     `json:"uVolume"`
	Aopen   float64 `json:"open"`
	Ahigh   float64 `json:"high"`
	Alow    float64 `json:"low"`
	Aclose  float64 `json:"close"`
	Avolume int     `json:"volume"`
}

// PriceHistory models historical prices for a security
type PriceHistory struct {
	Symbol string
	Prices []Prices
}

// Dividend type models a dividend event
type Dividend struct {
	DecDate     string `json:"declaredDate"`
	ExDate      string `json:"exDate"`
	RecDate     string `json:"recordDate"`
	PayDate     string `json:"paymentDate"`
	Amount      string `json:"amount"`
	Flag        string `json:"flag"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
}

// DividendHistory models dividends for a security
type DividendHistory struct {
	Symbol    string
	Dividends []Dividend
}

// IsEmpty checks whether dividend d is empty
func (dh DividendHistory) IsEmpty() bool {
	return len(dh.Dividends) == 0
}

// Split type models a security split event
type Split struct {
	DecDate     string  `json:"declaredDate"`
	ExDate      string  `json:"exDate"`
	Ratio       float64 `json:"ratio"`
	ToFactor    float64 `json:"toFactor"`
	FromFactor  float64 `json:"fromFactor"`
	Description string  `json:"description"`
}

// SplitHistory models splits for a security
type SplitHistory struct {
	Symbol string
	Splits []Split
}

// IsEmpty checks whether split s is empty
func (sh SplitHistory) IsEmpty() bool {
	return len(sh.Splits) == 0
}
