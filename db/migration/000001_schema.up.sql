CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS companies (
    id       uuid NOT NULL DEFAULT gen_random_uuid(),
    name     TEXT NOT NULL,
    cik      int NOT NULL UNIQUE,
    website  TEXT,
    industry TEXT,
    sector   TEXT,
    ceo      TEXT,
    state    VARCHAR,
    city     TEXT,
    zip      VARCHAR,
    PRIMARY KEY (id)
    -- maybe add datetime added
    -- maybe add datetime ended
);

CREATE TABLE IF NOT EXISTS stocks (
    id      uuid NOT NULL DEFAULT gen_random_uuid(),
    comid   VARCHAR REFERENCES companies(id),
    symbol  VARCHAR NOT NULL UNIQUE,
    name    TEXT NOT NULL,
    sectype VARCHAR NOT NULL,
    active  BOOLEAN NOT NULL DEFAULT TRUE,
    curreny VARCHAR DEFAULT 'USD',
    country VARCHAR DEFAULT 'US'
);

CREATE TABLE IF NOT EXISTS prices (
    time    TIMESTAMP,
    secid   VARCHAR REFERENCES stocks(id),
    open    DOUBLE PRECISION,
    close   DOUBLE PRECISION,
    high    DOUBLE PRECISION,
    low     DOUBLE PRECISION,
    volume  int,
    uopen   DOUBLE PRECISION,
    uclose  DOUBLE PRECISION,
    uhigh   DOUBLE PRECISION,
    ulow    DOUBLE PRECISION,
    uvolume int,
    PRIMARY KEY(time, secid)
);

CREATE INDEX ON prices (time DESC, secid);

SELECT create_hypertable(
    'prices', 
    'time', 
    create_default_indexes => FALSE,
    chunk_time_interval => INTERVAL '1 day'
);

CREATE TABLE IF NOT EXISTS dividends (
    symbol    VARCHAR REFERENCES stocks(id),
    exDate    TIMESTAMP,
    decDate   TIMESTAMP,
    recDate   TIMESTAMP,
    payDate   TIMESTAMP,
    amount    DOUBLE PRECISION,
    flag      TEXT,
    currency  VARCHAR,
    frequency VARCHAR,
    PRIMARY KEY(symbol, exDate)
);

CREATE TABLE IF NOT EXISTS splits (
    symbol     VARCHAR REFERENCES stocks(id),
    exDate     TIMESTAMP,
    decDate    TIMESTAMP,
    ratio      NUMERIC,
    toFactor   NUMERIC,
    fromFactor NUMERIC,
    PRIMARY KEY(symbol, exDate)
);