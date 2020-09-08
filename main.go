package main

import (
	"defcor/app"
	"log"
	"os"
)

const lookback = "5y"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.LUTC)

	myapp, err := app.StartApplication()
	if err != nil {
		panic(err)
	}
	defer myapp.End()

	// symbs, err := myapp.DB.Symbols()
	// if err != nil {
	// 	panic(err)
	// }
	// rest := restOfStocks("AAT", symbs)
	// if err := myapp.Seed(rest, lookback); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(rest)
	if err := myapp.CompleteDividends("AAT", "5y"); err != nil {
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
