package db

import (
	"context"
	"defcor/iex"
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
func (c *Conn) InsertStock(s iex.Stock) (int, error) {
	sql := `INSERT INTO stocks (
		symbol, name, date_added, active, sectype, iexid, figi, currency, region, cik
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING secid`

	var id int
	if err := c.c.QueryRow(context.Background(), sql,
		s.Symbol, s.Name, s.Date, s.IsActive, s.Type, s.IexID, s.Figi, s.Currency, s.Region, s.Cik,
	).Scan(&id); err != nil {
		return -1, err // -1 is an invalid id
	}
	return id, nil
}

// Symbols returns a list of all symbols in the stocks table
func (c *Conn) Symbols() ([]string, error) {
	sql := `SELECT symbol FROM stocks`
	rows, err := c.c.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var syms []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		syms = append(syms, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return syms, nil
}
