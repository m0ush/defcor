package db

import (
	"context"
	"defcor/api"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

// Conn type stores a postgres database connectino
type Conn struct {
	c *pgx.Conn
}

// CreateConn creates a postgres connection struct
func CreateConn() (*Conn, error) {
	c, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &Conn{c}, nil
}

// Close ends a postgres connection
func (c *Conn) Close() error {
	return c.c.Close(context.Background())
}

// InsertStock inserts a stock record into a database
func (c *Conn) InsertStock(s api.Stock) error {
	sqlStatement := `
	INSERT INTO stocks (
		symbol, 
		name, 
		date_added, 
		active, 
		sectype,
		iexid,
		figi,
		currency,
		region,
		cik
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return nil
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		sqlStatement,
		s.Symbol,
		s.Name,
		s.Date,
		s.IsActive,
		s.Type,
		s.IexID,
		s.Figi,
		s.Currency,
		s.Region,
		s.Cik,
	)
	if err != nil {
		return nil
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

// SeedStocks inserts all stocks into the db
func (c *Conn) SeedStocks() error {
	stks, err := api.Securities()
	if err != nil {
		panic(err)
	}
	for _, s := range stks {
		fmt.Printf("Working on... %+v\n", s)
		if err := c.InsertStock(s); err != nil {
			return err
		}
	}
	return nil
}
