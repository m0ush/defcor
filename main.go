package main

import (
	"defcor/iex"
	"fmt"
)

func main() {
	// conn, err := db.CreateConn()
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()
	// stk := iex.Stock{
	// 	Symbol:   "AACG",
	// 	Name:     "vprioyateCtedoA GaDR AioSrT",
	// 	Date:     "2020-08-24",
	// 	Type:     "cs",
	// 	IexID:    "IEX_44595A4C53392D52",
	// 	Region:   "US",
	// 	Currency: "USD",
	// 	IsActive: true,
	// 	Figi:     "0BB02G06SV3P",
	// 	Cik:      39194,
	// }
	// id, err := conn.InsertStock(&stk)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(id)

	// client := iex.NewClient()
	// stocks, err := client.AllStocks()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, s := range stocks {
	// 	fmt.Println(s.Name)
	// }

	// client := iex.NewClient()
	// divs, err := client.Dividends("aapl", "2y")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(divs)

	// client := iex.NewClient()
	// splts, err := client.Splits("RNVA", "1y")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(splts)

	client := iex.NewClient()
	prices, err := client.Prices("AAPL", "1m")
	if err != nil {
		panic(err)
	}
	fmt.Println(prices)
}
