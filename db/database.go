package db

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
)

const baseURL = "https://sandbox.iexapis.com"

// Profile type models a company profile endpoint
type Profile struct {
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

// Company enhances data in Profile with SEC CIK code
type Company struct {
	Symbol      string
	CompanyName string
	Industry    string
	Website     string
	CEO         string
	Sector      string
	State       string
	City        string
	Zip         string
	Country     string
	CIK         int
}

// Stock type models a stock from a company endpoint
type Stock struct {
	Symbol      string
	CompanyName string
	IssueType   string
	Country     string
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
	// Register UUID data type for connection
	c.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: &pgtypeuuid.UUID{},
		Name:  "uuid",
		OID:   pgtype.UUIDOID,
	})
	return &Conn{c}, nil
}

// Disconnect ends a postgres connection
func (c *Conn) Disconnect() error {
	return c.c.Close(context.Background())
}

// NewProfile gets IEX data
func NewProfile(symbol string) (*Profile, error) {
	base, err := url.Parse(baseURL)
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

	var cp Profile
	if err := json.NewDecoder(resp.Body).Decode(&cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

// Company returns Company data
func (cp Profile) Company() Company {
	return Company{
		Symbol:      cp.Symbol,
		CompanyName: cp.CompanyName,
		Industry:    cp.Industry,
		Website:     cp.Website,
		CEO:         cp.CEO,
		Sector:      cp.Sector,
		State:       cp.State,
		City:        cp.City,
		Zip:         cp.Zip,
		Country:     cp.Country,
		CIK:         ccm.Find(cp.Symbol),
	}
}

// Stock returns Stock data
func (cp Profile) Stock() Stock {
	return Stock{
		Symbol:      cp.Symbol,
		CompanyName: cp.CompanyName,
		IssueType:   cp.IssueType,
		Country:     cp.Country,
	}
}

// InsertCompany inserts a company into a database
func (c *Conn) InsertCompany(co Company) error {
	sqlStatement := `
	INSERT INTO companies (name, cik, website, industry, sector, ceo, state, city, zip)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return nil
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), sqlStatement)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// InsertStock inserts a stock record into a database
func (c *Conn) InsertStock(stk Stock) error {
	sqlStatement := `
	INSERT INTO stocks (symbol, name, sectype)
	VALUES ($1, $2, $3)
	RETURNING id`

	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return nil
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), sqlStatement)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
