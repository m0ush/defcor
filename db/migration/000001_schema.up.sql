CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE IF NOT EXISTS stocks (
	secid serial,
	symbol varchar(6) NOT NULL UNIQUE,
	name varchar(30) NOT NULL,
	date_added date DEFAULT today (),
    active boolean DEFAULT TRUE,
    sectype varchar(3),
    iexid char(20),
    figi char(12),
    curreny char(3),
    region char(2),
    cik integer,
    CONSTRAINT pk_stocks PRIMARY KEY secid
);

CREATE TABLE IF NOT EXISTS prices (
	time timestamp,
	secid integer,
	open numeric(12, 6),
	close numeric(12, 6),
	high numeric(12, 6),
	low numeric(12, 6),
	volume integer,
	uopen numeric(12, 6),
	uclose numeric(12, 6),
	uhigh numeric(12, 6),
	ulow numeric(12, 6),
	uvolume integer,
	CONSTRAINT pk_prices PRIMARY KEY (time, secid),
    CONSTRAINT fk_stocks_prices FOREIGN KEY secid REFERENCES stocks (secid)
);

CREATE TABLE IF NOT EXISTS dividends (
	secid integer,
	exDate timestamp,
	decDate timestamp,
	recDate timestamp,
	payDate timestamp,
	amount numeric(5, 2),
	flag varchar(120),
	currency varchar(4),
	frequency varchar(20),
	CONSTRAINT pk_dividends PRIMARY KEY (secid, exDate),
    CONSTRAINT fk_stocks_dividends FOREIGN KEY secid REFERENCES stocks (secid)
);

CREATE TABLE IF NOT EXISTS splits (
	secid integer,
	exDate timestamp,
	decDate timestamp,
	ratio numeric(9, 6),
	toFactor numeric(4, 2),
	fromFactor numeric(4, 2),
	CONSTRAINT pk_splits PRIMARY KEY (secid, exDate),
    CONSTRAINT fk_stocks_splits FOREIGN KEY secid REFERENCES stocks (secid)
);

CREATE INDEX ON prices (time DESC, secid);

SELECT
	create_hypertable ('prices',
		'time',
		create_default_indexes => FALSE,
		chunk_time_interval => interval '1 day');