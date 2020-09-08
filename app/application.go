package app

import (
	"context"
	"defcor/db"
	"defcor/iex"
)

// Application combines db with api
type Application struct {
	DB  *db.Conn
	api *iex.APIConnection
}

// StartApplication creates an app
func StartApplication() (*Application, error) {
	conn, err := db.CreateConn()
	if err != nil {
		return nil, err
	}
	api := iex.NewAPIConnection()
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
func (app *Application) CompletePrices(symbol, lookback string) error {
	securityPrices, err := app.api.Prices(context.Background(), symbol, lookback)
	if err != nil {
		return err
	}
	if err := app.DB.InsertPriceHistory(securityPrices); err != nil {
		return err
	}
	return nil
}

// CompleteDividends fetches dividends and inserts them into the database
func (app *Application) CompleteDividends(symbol, lookback string) error {
	securityDivs, err := app.api.Dividends(context.Background(), symbol, lookback)
	if err != nil {
		return err
	}
	if err := app.DB.InsertDividendHistory(securityDivs); err != nil {
		return err
	}
	return nil
}

// CompleteSplits fetches splits and inserts them into the database
func (app *Application) CompleteSplits(symbol, lookback string) error {
	securitySplits, err := app.api.Splits(context.Background(), symbol, lookback)
	if err != nil {
		return err
	}
	if err := app.DB.InsertSplitHistory(securitySplits); err != nil {
		return err
	}
	return nil
}

// Seed populates the application database
func (app *Application) Seed(symbols []string, lookback string) error {
	for _, symb := range symbols {
		if err := app.CompletePrices(symb, lookback); err != nil {
			return err
		}
		if err := app.CompleteDividends(symb, lookback); err != nil {
			return err
		}
		if err := app.CompleteSplits(symb, lookback); err != nil {
			return err
		}
	}
	return nil
}

// NewStocks add only new securities to the stocks table
func (app *Application) NewStocks() ([]iex.Stock, error) {
	exists, err := app.DB.Stocks()
	if err != nil {
		return nil, err
	}
	stocks, err := app.api.AllStocks(context.Background())
	if err != nil {
		return nil, err
	}
	newStocks := setDifference(exists, stocks)
	return newStocks, nil
}

func setDifference(existing, recent []iex.Stock) []iex.Stock {
	exists := make(map[iex.Stock]struct{}, len(existing))
	for _, x := range existing {
		exists[x] = struct{}{}
	}
	var diff []iex.Stock
	for _, x := range recent {
		if _, found := exists[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
