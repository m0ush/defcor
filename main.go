package main

import (
	"context"
	"defcor/db"
	"defcor/iex"
	"log"
	"os"
)

const lookback = "5y"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.LUTC)

	api := iex.NewAPIConnection()

	db, err := db.CreateConn()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()

	stocks, err := api.AllStocks(ctx)
	if err != nil {
		panic(err)
	}

	for _, stock := range stocks {
		log.Printf("Tracking: %s %s\n", stock.Symbol, stock.Name)
		_, err := db.InsertStock(stock)
		if err != nil {
			panic(err)
		}
	}

	symbols, err := db.Symbols()
	if err != nil {
		panic(err)
	}

	for _, symb := range symbols {
		sPrices, err := api.Prices(ctx, symb, lookback)
		if err != nil {
			panic(err)
		}
		log.Printf("Last Price: %s %v\n", sPrices.Symbol, sPrices.Prices[0])
		if err := db.InsertPriceHistory(sPrices); err != nil {
			panic(err)
		}

		sDivs, err := api.Dividends(ctx, symb, lookback)
		if err != nil {
			panic(err)
		}
		log.Printf("Last Dividend: %s %v\n", sDivs.Symbol, sDivs.Dividends[0])
		if err := db.InsertDividendHistory(sDivs); err != nil {
			panic(err)
		}

		sSplits, err := api.Splits(ctx, symb, lookback)
		if err != nil {
			panic(err)
		}
		log.Printf("Last Split: %s %v\n", sSplits.Symbol, sSplits.Splits[0])
		if err := db.InsertSplitHistory(sSplits); err != nil {
			panic(err)
		}
	}
	log.Println("Done")
}
