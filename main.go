package main

import (
	"defcor/db"
	"fmt"
)

func main() {
	aapl, err := db.Profile("AAPL")
	if err != nil {
		panic(err)
	}
	fmt.Println(aapl)
}
