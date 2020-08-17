package db

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/jackc/pgx/v4"
)

const baseURL = "https://sandbox.iexapis.com"

// Company type models a company profile endpoint
type Company struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	Industry    string `json:"industry"`
	Website     string `json:"website"`
	CEO         string `json:"CEO"`
	IssueType   string `json:"issueType"`
	Sector      string `json:"sector"`
	State       string `json:"state"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
}

// CompanyWithCik enhances data in Company with SEC CIK code
type CompanyWithCik struct {
	Company
	CIK int
}

// Stock type models a stock from a company endpoint
type Stock struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	IssueType   string `json:"issueType"`
	Country     string `json:"country"`
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

// Conn type stores a postgres database connectino
type Conn struct {
	c *pgx.Conn
}

// CreateConn creates a postgres connection struct
func CreateConn() (*Conn, error) {
	c, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Conn{c}, nil
}

// Disconnect ends a postgres connection
func (c *Conn) Disconnect() error {
	return c.c.Close(context.Background())
}

// Profile retrieves a company profile
func Profile(symbol string) (*CompanyWithCik, error) {
	base, err := url.Parse("https://sandbox.iexapis.com")
	if err != nil {
		return nil, err
	}
	relative, err := url.Parse(path.Join("stable", "stock", symbol, "company"))
	if err != nil {
		return nil, err
	}
	params := relative.Query()
	params.Set("token", os.Getenv("IEXCLOUD_TEST"))
	relative.RawQuery = params.Encode()

	u := base.ResolveReference(relative)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var company Company
	if err := json.NewDecoder(resp.Body).Decode(&company); err != nil {
		return nil, err
	}

	cik := ccm.Find(symbol)

	return &CompanyWithCik{company, cik}, nil
}

// InsertCompany inserts a company into a database
func (c *Conn) InsertCompany(co *Company) error {
	return nil
}
