CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE;

CREATE TABLE IF NOT EXISTS stocks (
	secid serial PRIMARY KEY,
	symbol varchar(6) NOT NULL UNIQUE,
	name varchar(50) NOT NULL,
	date_added date DEFAULT CURRENT_DATE,
    active boolean DEFAULT TRUE,
    sectype varchar(3),
    iexid char(20),
    figi char(12),
    currency char(3),
    region char(2),
    cik integer
);

CREATE TABLE IF NOT EXISTS prices (
	time timestamp,
	secid integer REFERENCES stocks (secid),
	uopen numeric(12, 6),
	uclose numeric(12, 6),
	uhigh numeric(12, 6),
	ulow numeric(12, 6),
	uvolume integer,
	aopen numeric(12, 6),
	aclose numeric(12, 6),
	ahigh numeric(12, 6),
	alow numeric(12, 6),
	avolume integer,
	PRIMARY KEY (time, secid)
);

CREATE TABLE IF NOT EXISTS dividends (
	secid integer REFERENCES stocks (secid),
	exDate timestamp,
	decDate timestamp,
	recDate timestamp,
	payDate timestamp,
	amount numeric(5, 2),
	flag varchar(120),
	currency varchar(4),
	frequency varchar(20),
    PRIMARY KEY (secid, exDate)
);

CREATE TABLE IF NOT EXISTS splits (
	secid integer REFERENCES stocks (secid),
	exDate timestamp,
	decDate timestamp,
	ratio numeric(10, 6),
	toFactor numeric(7, 2),
	fromFactor numeric(7, 2),
	PRIMARY KEY (secid, exDate)
);

CREATE INDEX ON prices (time DESC, secid);

SELECT
	create_hypertable ('prices',
		'time',
		create_default_indexes => FALSE,
		chunk_time_interval => interval '1 day');