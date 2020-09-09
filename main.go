package main

import (
	"defcor/db"
	"defcor/iex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// const lookback = "5y"

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.LUTC)

	f, err := os.Open("tester.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytedata, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var ds []iex.Dividend
	if err := json.Unmarshal(bytedata, &ds); err != nil {
		panic(err)
	}
	for _, d := range ds {
		fmt.Println(d)
	}
	conn, err := db.CreateConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if err := conn.TestDivInsert(233, ds); err != nil {
		panic(err)
	}
	if err := conn.TestDivDeletes(233); err != nil {
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
