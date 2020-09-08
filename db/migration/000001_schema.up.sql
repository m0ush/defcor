CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE;

CREATE TABLE IF NOT EXISTS stocks (
	secid serial PRIMARY KEY,
	symbol varchar(6) NOT NULL UNIQUE,
	name varchar(120) NOT NULL,
	date_added date,
    active boolean,
    sectype char(2),
    iexid char(20) UNIQUE,
    figi char(12) UNIQUE,
    currency char(3),
    region char(2),
    cik integer
);

-- Likely need to partition verticially by unadj/adj
CREATE TABLE IF NOT EXISTS prices (
	date date,
	secid integer REFERENCES stocks (secid),
	uopen numeric(10, 4),
	uclose numeric(10, 4),
	uhigh numeric(10, 4),
	ulow numeric(10, 4),
	uvolume integer,
	aopen numeric(10, 4),
	aclose numeric(10, 4),
	ahigh numeric(10, 4),
	alow numeric(10, 4),
	avolume integer,
	PRIMARY KEY (date, secid)
);

CREATE TABLE IF NOT EXISTS dividends (
	divid serial PRIMARY KEY,
	secid integer REFERENCES stocks (secid),
	decdate date,
	exdate date,
	recdate date,
	paydate date,
	amount numeric(8, 4),
	flag varchar(20),
	currency char(3),
	frequency varchar(20),
	description varchar(120)
);

CREATE TABLE IF NOT EXISTS splits (
	secid integer REFERENCES stocks (secid),
	decdate date,
	exdate date,
	ratio numeric(10, 6),
	tofactor numeric(7, 2),
	fromfactor numeric(7, 2),
	description varchar(120),
	PRIMARY KEY (secid, exdate)
);

CREATE INDEX ON prices (date DESC, secid);
CREATE INDEX ON dividends (secid, exdate DESC);

SELECT
	create_hypertable ('prices',
		'date',
		create_default_indexes => FALSE,
		chunk_time_interval => interval '1 day');