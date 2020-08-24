package main

import (
	"defcor/api"
	"defcor/db"
	"fmt"
)

func main() {
	conn, err := db.CreateConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	stk := api.Stock{
		Symbol:   "AACG",
		Name:     "vprioyateCtedoA GaDR AioSrT",
		Date:     "2020-08-24",
		Type:     "cs",
		IexID:    "IEX_44595A4C53392D52",
		Region:   "US",
		Currency: "USD",
		IsActive: true,
		Figi:     "0BB02G06SV3P",
		Cik:      39194,
	}
	id, err := conn.InsertStock(&stk)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
