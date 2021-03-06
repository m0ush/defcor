package db

import (
	"context"
	"database/sql"
	"defcor/iex"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

var tfmt = "2006-01-02"

// Conn type stores a postgres database connection
type Conn struct {
	c *pgx.Conn
}

// CreateConn creates a postgres connection struct
func CreateConn(dburl string) (*Conn, error) {
	c, err := pgx.Connect(context.Background(), dburl)
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
		symbol, name, date_added, sectype, iexid, figi, currency, region, cik
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING secid`
	today := time.Now().Format(tfmt)

	var secid int
	if err := c.c.QueryRow(context.Background(), sql,
		s.Symbol, s.Name, today, s.Type, s.IexID, s.Figi, s.Curr, s.Region, s.Cik,
	).Scan(&secid); err != nil {
		return -1, err // -1 is an invalid secid
	}
	return secid, nil
}

// InsertStocks adds a slice of stocks into the stocks database
func (c *Conn) InsertStocks(stks []iex.Stock) error {
	for _, stk := range stks {
		secid, err := c.InsertStock(stk)
		if err != nil {
			return err
		}
		log.Printf("secid: %d(%s, %s)", secid, stk.Symbol, stk.Name)
	}
	return nil
}

// Symbols returns a list of all symbols in the stocks table
func (c *Conn) Symbols() ([]string, error) {
	sql := `SELECT symbol FROM stocks WHERE date_inactive IS NULL`
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

// Stocks returns the entire stock db
func (c *Conn) Stocks() ([]iex.Stock, error) {
	sql := `SELECT symbol, name, sectype, iexid, figi, currency, region, cik FROM stocks WHERE date_inactive IS NULL`
	rows, err := c.c.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stks []iex.Stock
	for rows.Next() {
		var s iex.Stock
		if err := rows.Scan(
			&s.Symbol, &s.Name, &s.Type, &s.IexID, &s.Figi, &s.Curr, &s.Region, &s.Cik,
		); err != nil {
			return nil, err
		}
		stks = append(stks, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stks, nil
}

// FindSecurityID grabs the secid from the stocks table
func (c *Conn) FindSecurityID(symbol string) (int, error) {
	var secid int
	sql := `SELECT secid FROM stocks WHERE symbol=$1`
	row := c.c.QueryRow(context.Background(), sql, symbol)
	err := row.Scan(&secid)
	if err != nil {
		return -1, err
	}
	return secid, nil
}

// InsertPriceHistory inserts a stock's historical price data into the price table
func (c *Conn) InsertPriceHistory(ph *iex.PriceHistory) error {
	sql := `INSERT INTO prices(
		date, secid, uopen, uclose, uhigh, ulow, uvolume, aopen, aclose, ahigh, alow, avolume
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	secid, err := c.FindSecurityID(ph.Symbol)
	if err != nil {
		return err
	}
	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return err
	}
	// no-op on successful tx commit
	defer tx.Rollback(context.Background())
	for _, p := range ph.Prices {
		_, err = tx.Exec(context.Background(), sql,
			p.Date, secid, p.Uopen, p.Uclose, p.Uhigh, p.Ulow, p.Uvolume, p.Aopen, p.Aclose, p.Ahigh, p.Alow, p.Avolume,
		)
		if err != nil {
			return fmt.Errorf("insertion error: %s(price:%v)", ph.Symbol, p)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// InsertDividendHistory inserts a stock's historical dividends into the dividends table
func (c *Conn) InsertDividendHistory(dh *iex.DividendHistory) error {
	if dh.IsEmpty() {
		return nil
	}
	sql := `INSERT INTO dividends
		(secid, decdate, exdate, recdate, paydate, amount, flag, currency, frequency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	secid, err := c.FindSecurityID(dh.Symbol)
	if err != nil {
		return err
	}
	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return err
	}
	// no-op on successful tx commit
	defer tx.Rollback(context.Background())
	for _, d := range dh.Dividends {
		_, err = tx.Exec(context.Background(), sql,
			secid, d.DecDate, d.ExDate, d.RecDate, d.PayDate, d.Amount, d.Flag, d.Curr, d.Freq,
		)
		if err != nil {
			return fmt.Errorf("insertion error: %s(dividend:%v)", dh.Symbol, d)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// InsertSplitHistory inserts a stock's historical stock splits into the splits table
func (c *Conn) InsertSplitHistory(sh *iex.SplitHistory) error {
	if sh.IsEmpty() {
		return nil
	}
	sql := `INSERT INTO splits(secid, decdate, exdate, tofactor, fromfactor)
		VALUES ($1, $2, $3, $4, $5)`

	secid, err := c.FindSecurityID(sh.Symbol)
	if err != nil {
		return err
	}
	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return err
	}
	// no-op on successful tx commit
	defer tx.Rollback(context.Background())
	for _, s := range sh.Splits {
		_, err = tx.Exec(context.Background(), sql,
			secid, s.DecDate, s.ExDate, s.ToFactor, s.FromFactor,
		)
		if err != nil {
			return fmt.Errorf("insertion error: %s(split:%v)", sh.Symbol, s)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func nullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// TestDivInsert tests to make sure dividends are inserted well
func (c *Conn) TestDivInsert(i int, ds []iex.Dividend) error {
	sql := `INSERT INTO dividends
		(secid, decdate, exdate, recdate, paydate, amount, flag, currency, frequency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	tx, err := c.c.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, d := range ds {
		_, err = tx.Exec(context.Background(), sql,
			i, d.DecDate, d.ExDate, d.RecDate, d.PayDate, d.Amount, d.Flag, d.Curr, d.Freq,
		)
		if err != nil {
			return fmt.Errorf("Error on insert %v", d)
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// TestDivDeletes deletes a dividend record where secid = i
func (c *Conn) TestDivDeletes(i int) error {
	sql := `DELETE FROM dividends WHERE secid=$1`
	_, err := c.c.Exec(context.Background(), sql, i)
	return err
}
