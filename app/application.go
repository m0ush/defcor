package app

import (
	"context"
	"defcor/db"
	"defcor/iex"
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
	getSome(existing, refreshed)
	return nil
}

func getSome(A, B iex.StockGroup) {
	iex.FormatOutput(A, B)
}
