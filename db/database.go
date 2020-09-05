package db

import (
	"context"
	"database/sql"
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

	var secid int
	if err := c.c.QueryRow(context.Background(), sql,
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
	).Scan(&secid); err != nil {
		return -1, err // -1 is an invalid secid
	}
	return secid, nil
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
			p.Date,
			secid,
			p.Uopen,
			p.Uclose,
			p.Uhigh,
			p.Ulow,
			p.Uvolume,
			p.Aopen,
			p.Aclose,
			p.Ahigh,
			p.Alow,
			p.Avolume,
		)
		if err != nil {
			return err
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
	sql := `INSERT INTO dividends
		(secid, decdate, exdate, recdate, paydate, amount, flag, currency, frequency, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// Do not insert if there's no dividend history
	if dh.IsEmpty() {
		return nil
	}

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
			secid,
			d.DecDate,
			d.ExDate,
			d.RecDate,
			nullString(d.PayDate),
			nullString(d.Amount),
			nullString(d.Flag),
			nullString(d.Currency),
			nullString(d.Frequency),
			nullString(d.Description),
		)
		if err != nil {
			return err
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
	sql := `INSERT INTO splits(
		secid, decdate, exdate, ratio, tofactor, fromfactor, description
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Don't insert if there's no split history
	if sh.IsEmpty() {
		return nil
	}

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
			secid,
			s.DecDate,
			s.ExDate,
			s.Ratio,
			s.ToFactor,
			s.FromFactor,
			nullString(s.Description),
		)
		if err != nil {
			return err
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
