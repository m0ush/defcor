package main

import (
	"defcor/db"
	"fmt"
)

func main() {
	conn, err := db.CreateConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := conn.SeedStocks(); err != nil {
		fmt.Println(err)
	}
}
