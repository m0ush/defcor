package main

import (
	"defcor/app"
	"defcor/iex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var environment = app.Environment{
	Host:     "sandbox.iexapis.com",
	APIKey:   os.Getenv("IEXCLOUD_TEST"),
	Lookback: "5y",
	Duration: 75 * time.Millisecond,
	DbURL:    os.Getenv("DATABASE_URL_TEST"),
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.LUTC)

	myapp, err := app.StartApplication(environment)
	if err != nil {
		panic(err)
	}
	defer myapp.End()
	// if err := myapp.RefreshStocks(); err != nil {
	// 	panic(err)
	// }
	symbols, err := myapp.DB.Symbols()
	if err != nil {
		panic(err)
	}
	revised := restOfStocks("MCRI", symbols)
	if err := myapp.Seed(revised); err != nil {
		panic(err)
	}
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
