package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
)

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

type Stock struct {
	Symbol      string `json:"symbol"`
	CompanyName string `json:"companyName"`
	IssueType   string `json:"issueType"`
	Country     string `json:"country"`
}

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

type Split []struct {
	ExDate       string  `json:"exDate"`
	DeclaredDate string  `json:"declaredDate"`
	Ratio        float64 `json:"ratio"`
	ToFactor     int     `json:"toFactor"`
	FromFactor   int     `json:"fromFactor"`
	Description  string  `json:"description"`
}

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

func makeConn() *pgx.Conn {
	conn, _ := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	return conn
}
