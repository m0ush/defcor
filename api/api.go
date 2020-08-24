package api

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
	Cik      int    `json:"cik"`
}

// Dividend type models a dividend series endpoint
type Dividend []struct {
	Symbol       string  `json:"symbol"`
	ExDate       string  `json:"exDate"`
	PaymentDate  string  `json:"paymentDate"`
	RecordDate   string  `json:"recordDate"`
	DeclaredDate string  `json:"declaredDate"`
	Amount       float64 `json:"amount"`
	Flag         string  `json:"flag"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
	Frequency    string  `json:"frequency"`
}

// Split type models a split series endpoint
// Note: does not include symbol by default
type Split []struct {
	ExDate       string  `json:"exDate"`
	DeclaredDate string  `json:"declaredDate"`
	Ratio        float64 `json:"ratio"`
	ToFactor     int     `json:"toFactor"`
	FromFactor   int     `json:"fromFactor"`
	Description  string  `json:"description"`
}

// PriceHistory type models a price history series endpoint
type PriceHistory []struct {
	Date    string  `json:"date"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
	Volume  int     `json:"volume"`
	UOpen   float64 `json:"uOpen"`
	UHigh   float64 `json:"uHigh"`
	ULow    float64 `json:"uLow"`
	UClose  float64 `json:"uClose"`
	UVolume int     `json:"uVolume"`
}
