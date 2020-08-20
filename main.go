package main

import (
	"defcor/db"
	"fmt"
)

func main() {
	aapl, err := db.NewProfile("AAPL")
	if err != nil {
		panic(err)
	}
	co := aapl.Company()
	stk := aapl.Stock()

	fmt.Println(co)
	fmt.Println(stk)
}
