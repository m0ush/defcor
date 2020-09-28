package main

import (
	"defcor/app"
	"defcor/db"
	"defcor/iex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var environment = app.Environment{
	Host:     "cloud.iexapis.com",
	APIKey:   os.Getenv("IEXCLOUD_SECRET"),
	Lookback: "5d",
	Duration: 60 * time.Millisecond,
	DbURL:    os.Getenv("DATABASE_URL_PROD"),
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.LUTC)
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	myapp, err := app.Start(environment)
	if err != nil {
		return fmt.Errorf("setting up app: %w", err)
	}
	defer myapp.End()
	if err := myapp.RefreshStocks(); err != nil {
		panic(err)
	}
	symbols, err := myapp.DB.Symbols()
	if err != nil {
		return fmt.Errorf("getting symbols: %w", err)
	}
	// revised := restOfStocks("WUBA", symbols)
	if err := myapp.Seed(symbols); err != nil {
		return fmt.Errorf("seeding problem: %w", err)
	}
	return nil
}

func restOfStocks(element string, data []string) []string {
	var x int
	for k, v := range data {
		if v == element {
			x = k
		}
	}
	return data[x+1:]
}

func readtesterfile() ([]iex.Dividend, error) {
	f, err := os.Open("tester.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytedata, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var ds []iex.Dividend
	if err := json.Unmarshal(bytedata, &ds); err != nil {
		return nil, err
	}
	return ds, nil
}

func customupload(filename string) ([]iex.PriceHistory, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytedata, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var this map[string][]iex.Prices
	if err := json.Unmarshal(bytedata, &this); err != nil {
		return nil, err
	}
	var allPrices []iex.PriceHistory
	for k, entry := range this {
		allPrices = append(allPrices, iex.PriceHistory{
			Symbol: k,
			Prices: entry,
		})
	}
	return allPrices, nil
}

func insertPriceUpload(filename string) error {
	conn, err := db.CreateConn(environment.DbURL)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer conn.Close()
	this, err := customupload(filename)
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}
	for _, entry := range this {
		conn.InsertPriceHistory(&entry)
	}
	return nil
}
