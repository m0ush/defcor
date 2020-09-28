package iex

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// NullString mirrors a sql.NullString
type NullString struct{ sql.NullString }

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

// NullNumber wraps a sql.NullFloat64 around a json.NullNumber
type NullNumber struct{ sql.NullFloat64 }

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

// NullInt64 mirrors a sql.NullInt64
type NullInt64 struct{ sql.NullInt64 }

// UnmarshalJSON checks to see if value is an integer or null
// func (ni *NullInt64) UnmarshalJSON(data []byte) error {
// 	if bytes.Equal(data, []byte(`null`)) {
// 		ni.Valid = false
// 		return nil
// 	}
// 	var x int64
// 	if err := json.Unmarshal(data, &x); err != nil {
// 		return err
// 	}
// 	ni.Valid = true
// 	ni.Int64 = x
// 	return nil
// }

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

// IncomeHistory models historical income statements for a security
type IncomeHistory struct {
	Symbol string `json:"symbol"`
	Income []struct {
		ReportDate             string `json:"reportDate"`
		FiscalDate             string `json:"fiscalDate"`
		Currency               string `json:"currency"`
		TotalRevenue           int64  `json:"totalRevenue"`
		CostOfRevenue          int64  `json:"costOfRevenue"`
		GrossProfit            int64  `json:"grossProfit"`
		ResearchAndDevelopment int64  `json:"researchAndDevelopment"`
		SellingGeneralAndAdmin int64  `json:"sellingGeneralAndAdmin"`
		OperatingExpense       int64  `json:"operatingExpense"`
		OperatingIncome        int64  `json:"operatingIncome"`
		OtherIncomeExpenseNet  int64  `json:"otherIncomeExpenseNet"`
		Ebit                   int64  `json:"ebit"`
		InterestIncome         int64  `json:"interestIncome"`
		PretaxIncome           int64  `json:"pretaxIncome"`
		IncomeTax              int64  `json:"incomeTax"`
		MinorityInterest       int64  `json:"minorityInterest"`
		NetIncome              int64  `json:"netIncome"`
		NetIncomeBasic         int64  `json:"netIncomeBasic"`
	} `json:"income"`
}

// BalanceHistory models the balance sheet history for a security
type BalanceHistory struct {
	Symbol       string `json:"symbol"`
	Balancesheet []struct {
		ReportDate              string    `json:"reportDate"`
		FiscalDate              string    `json:"fiscalDate"`
		Currency                string    `json:"currency"`
		CurrentCash             int64     `json:"currentCash"`
		ShortTermInvestments    int64     `json:"shortTermInvestments"`
		Receivables             int64     `json:"receivables"`
		Inventory               int64     `json:"inventory"`
		OtherCurrentAssets      int64     `json:"otherCurrentAssets"`
		CurrentAssets           int64     `json:"currentAssets"`
		LongTermInvestments     int64     `json:"longTermInvestments"`
		PropertyPlantEquipment  int64     `json:"propertyPlantEquipment"`
		Goodwill                NullInt64 `json:"goodwill"`
		IntangibleAssets        NullInt64 `json:"intangibleAssets"`
		OtherAssets             int64     `json:"otherAssets"`
		TotalAssets             int64     `json:"totalAssets"`
		AccountsPayable         int64     `json:"accountsPayable"`
		CurrentLongTermDebt     int64     `json:"currentLongTermDebt"`
		OtherCurrentLiabilities int64     `json:"otherCurrentLiabilities"`
		TotalCurrentLiabilities int64     `json:"totalCurrentLiabilities"`
		LongTermDebt            int64     `json:"longTermDebt"`
		OtherLiabilities        int64     `json:"otherLiabilities"`
		MinorityInterest        int64     `json:"minorityInterest"`
		TotalLiabilities        int64     `json:"totalLiabilities"`
		CommonStock             int64     `json:"commonStock"`
		RetainedEarnings        int64     `json:"retainedEarnings"`
		TreasuryStock           NullInt64 `json:"treasuryStock"`
		CapitalSurplus          NullInt64 `json:"capitalSurplus"`
		ShareholderEquity       int64     `json:"shareholderEquity"`
		NetTangibleAssets       int64     `json:"netTangibleAssets"`
	} `json:"balancesheet"`
}

// CashFlowHistory models an iex Cashflow entry
type CashFlowHistory struct {
	Symbol   string `json:"symbol"`
	Cashflow []struct {
		ReportDate              string    `json:"reportDate"`
		FiscalDate              string    `json:"fiscalDate"`
		Currency                string    `json:"currency"`
		NetIncome               int64     `json:"netIncome"`
		Depreciation            int64     `json:"depreciation"`
		ChangesInReceivables    int64     `json:"changesInReceivables"`
		ChangesInInventories    int64     `json:"changesInInventories"`
		CashChange              int64     `json:"cashChange"`
		CashFlow                int64     `json:"cashFlow"`
		CapitalExpenditures     int64     `json:"capitalExpenditures"`
		Investments             int64     `json:"investments"`
		InvestingActivityOther  int64     `json:"investingActivityOther"`
		TotalInvestingCashFlows int64     `json:"totalInvestingCashFlows"`
		DividendsPaid           int64     `json:"dividendsPaid"`
		NetBorrowings           int64     `json:"netBorrowings"`
		OtherFinancingCashFlows int64     `json:"otherFinancingCashFlows"`
		CashFlowFinancing       int64     `json:"cashFlowFinancing"`
		ExchangeRateEffect      NullInt64 `json:"exchangeRateEffect"`
	} `json:"cashflow"`
}
