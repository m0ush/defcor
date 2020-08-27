package main

import (
	"defcor/db"
	"defcor/iex"
	"fmt"
)

func main() {
	conn, err := db.CreateConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := iex.NewClient()
	stocks, err := client.AllStocks()
	if err != nil {
		panic(err)
	}

	for _, s := range stocks {
		fmt.Printf("Stock:%v-`%v`;%v\n", s.Symbol, s.Cik, len(s.Cik))
		// if s.Cik == " " {
		// 	continue
		// }
		// id, err := conn.InsertStock(s)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(id)
	}
}
