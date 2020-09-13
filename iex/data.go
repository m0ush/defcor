package iex

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// NullString mirrors a sql.NullString
type NullString struct{ sql.NullString }

// NullNumber wraps a sql.NullFloat64 around a json.NullNumber
type NullNumber struct{ sql.NullFloat64 }

// UnmarshalJSON checks to see if a value is either null or empty
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || len(data) == 2 {
		ns.Valid = false
		return nil
	}
	if err := json.Unmarshal(data, &ns.String); err != nil {
		return err
	}
	ns.Valid = true
	return nil
}

// UnmarshalJSON checks to see if a value is either a json.Number or empty string
func (nn *NullNumber) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte(`""`)) {
		nn.Valid = false
		return nil
	}
	var x *json.Number
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	this, err := x.Float64()
	if err != nil {
		return err
	}
	nn.Valid = true
	nn.Float64 = this
	return nil
}

// Stock represents a single Stock
type Stock struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	IexID    string `json:"iexId"`
	Region   string `json:"region"`
	Curr     string `json:"currency"`
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
	DecDate NullString `json:"declaredDate"`
	ExDate  string     `json:"exDate"`
	RecDate NullString `json:"recordDate"`
	PayDate NullString `json:"paymentDate"`
	Amount  NullNumber `json:"amount"`
	Flag    NullString `json:"flag"`
	Curr    NullString `json:"currency"`
	Desc    NullString `json:"description"`
	Freq    NullString `json:"frequency"`
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
	DecDate    NullString `json:"declaredDate"`
	ExDate     string     `json:"exDate"`
	Ratio      float64    `json:"ratio"`
	ToFactor   float64    `json:"toFactor"`
	FromFactor float64    `json:"fromFactor"`
	Desc       NullString `json:"description"`
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
