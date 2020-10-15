package app

import (
	"context"
	"defcor/db"
	"defcor/iex"
	"fmt"
	"log"
	"time"
)

// Application combines db with api
type Application struct {
	DB  *db.Conn
	api *iex.APIConnection
}

// Environment outlines the environment
type Environment struct {
	Host     string
	APIKey   string
	Lookback string
	Duration time.Duration
	DbURL    string
}

// Start creates an app
func Start(env Environment) (*Application, error) {
	conn, err := db.CreateConn(env.DbURL)
	if err != nil {
		return nil, err
	}
	api := iex.NewAPIConnection(env.Host, env.APIKey, env.Lookback, env.Duration)
	return &Application{
		DB:  conn,
		api: api,
	}, nil
}

// End closes the database
func (app *Application) End() error {
	return app.DB.Close()
}

// CompletePrices fetches prices and inserts them into the database
func (app *Application) CompletePrices(symbol string) error {
	securityPrices, err := app.api.Prices(context.Background(), symbol)
	if err != nil {
		return err
	}
	if err := app.DB.InsertPriceHistory(securityPrices); err != nil {
		return err
	}
	return nil
}

// CompleteDividends fetches dividends and inserts them into the database
func (app *Application) CompleteDividends(symbol string) error {
	securityDivs, err := app.api.Dividends(context.Background(), symbol)
	if err != nil {
		return err
	}
	if err := app.DB.InsertDividendHistory(securityDivs); err != nil {
		return err
	}
	return nil
}

// CompleteSplits fetches splits and inserts them into the database
func (app *Application) CompleteSplits(symbol string) error {
	securitySplits, err := app.api.Splits(context.Background(), symbol)
	if err != nil {
		return err
	}
	if err := app.DB.InsertSplitHistory(securitySplits); err != nil {
		return err
	}
	return nil
}

// Seed populates the application database
func (app *Application) Seed(symbols []string) error {
	for _, symb := range symbols {
		log.Printf("working on %s...\n", symb)
		if err := app.CompletePrices(symb); err != nil {
			return err
		}
		log.Println("prices complete")
		if err := app.CompleteDividends(symb); err != nil {
			return err
		}
		log.Println("dividends complete")
		if err := app.CompleteSplits(symb); err != nil {
			return err
		}
		log.Println("splits complete")
	}
	log.Println("Done!")
	return nil
}

// RefreshStocks add only new securities to the stocks table
func (app *Application) RefreshStocks() error {
	existing, err := app.DB.Stocks()
	if err != nil {
		return err
	}
	refreshed, err := app.api.AllStocks(context.Background())
	if err != nil {
		return err
	}
	// newStocks := setDifference(existing, refreshed)
	// if err := app.DB.InsertStocks(newStocks); err != nil {
	// 	return err
	// }
	resolver := NewStockResolver(existing, refreshed)
	resolver.DifferencesFormatted()
	return nil
}

func setDifference(A, B []iex.Stock) []iex.Stock {
	SetA := make(map[iex.Stock]struct{}, len(A))
	for _, x := range A {
		SetA[x] = struct{}{}
	}
	var diff []iex.Stock
	for _, x := range B {
		if _, found := SetA[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// StockResolver is used to build sets for comparison between two Stock slices
type StockResolver struct {
	A []iex.Stock
	B []iex.Stock
}

// NewStockResolver creates a stock resolver
func NewStockResolver(existing, new []iex.Stock) StockResolver {
	return StockResolver{
		A: existing,
		B: new,
	}
}

// Differences returns the set differences of old and new stock data
func (sr StockResolver) Differences() ([]iex.Stock, []iex.Stock) {
	setA := make(map[iex.Stock]struct{}, len(sr.A))
	for _, x := range sr.A {
		setA[x] = struct{}{}
	}
	setB := make(map[iex.Stock]struct{}, len(sr.B))
	for _, x := range sr.B {
		setB[x] = struct{}{}
	}
	var ANotB []iex.Stock
	for _, x := range sr.A {
		if _, found := setB[x]; !found {
			ANotB = append(ANotB, x)
		}
	}
	var BNotA []iex.Stock
	for _, x := range sr.B {
		if _, found := setA[x]; !found {
			BNotA = append(BNotA, x)
		}
	}
	return ANotB, BNotA
}

// DifferencesFormatted provides a pretty printed version of the differences
func (sr StockResolver) DifferencesFormatted() {
	updates, a, b := sr.Reconcile()
	fmt.Println("\nUpdates:")
	for k, v := range updates {
		fmt.Println("\nreconcile...")
		fmt.Println("\tprev:", k)
		fmt.Println("\tcurr:", v)
	}
	fmt.Println("\nEnds:")
	for i, s := range a {
		fmt.Printf("%2d: %v\n", i, s)
	}
	fmt.Println("\nAdds:")
	for i, s := range b {
		fmt.Printf("%2d: %v\n", i, s)
	}
}

// Reconcile finds securities that need to be updated, added, or set to inactive
func (sr StockResolver) Reconcile() (map[iex.Stock]iex.Stock, []iex.Stock, []iex.Stock) {
	a, b := sr.Differences()
	mapA := make(map[string]int)
	for i, s := range a {
		mapA[s.IexID] = i
	}
	mapB := make(map[string]int)
	for i, s := range b {
		mapB[s.IexID] = i
	}
	var aRev []iex.Stock
	var tempB []iex.Stock
	updates := make(map[iex.Stock]iex.Stock)
	for k, idxa := range mapA {
		if idxb, ok := mapB[k]; ok {
			updates[a[idxa]] = b[idxb]
			tempB = append(tempB, b[idxb])
		} else {
			aRev = append(aRev, a[idxa])
		}
	}
	setBprime := make(map[iex.Stock]struct{}, len(tempB))
	for _, x := range tempB {
		setBprime[x] = struct{}{}
	}
	var BNotBprime []iex.Stock
	for _, x := range b {
		if _, found := setBprime[x]; !found {
			BNotBprime = append(BNotBprime, x)
		}
	}
	return updates, aRev, BNotBprime
}
