CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE companies (
    id       VARCHAR PRIMARY KEY,
    name     TEXT NOT NULL,
    cik      int NOT NULL UNIQUE,
    website  TEXT,
    industry TEXT,
    sector   TEXT,
    ceo      TEXT,
    state    VARCHAR,
    city     TEXT,
    zip      VARCHAR
);

CREATE TABLE stocks (
    id      VARCHAR PRIMARY KEY,
    comid   VARCHAR REFERENCES companies(id),
    symbol  VARCHAR NOT NULL UNIQUE,
    name    TEXT NOT NULL,
    sectype VARCHAR NOT NULL,
    active  BOOLEAN NOT NULL DEFAULT TRUE,
    curreny VARCHAR DEFAULT 'USD',
    country VARCHAR DEFAULT 'US'
);

CREATE TABLE prices (
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

CREATE TABLE dividends (
    id        VARCHAR PRIMARY KEY,
    symbol    VARCHAR REFERENCES stocks(id),
    exDate    TIMESTAMP,
    decDate   TIMESTAMP,
    recDate   TIMESTAMP,
    payDate   TIMESTAMP,
    amount    DOUBLE PRECISION,
    flag      TEXT,
    currency  VARCHAR,
    frequency VARCHAR
);

CREATE TABLE splits (
    id         VARCHAR PRIMARY KEY,
    symbol     VARCHAR REFERENCES stocks(id),
    exDate     TIMESTAMP,
    decDate    TIMESTAMP,
    ratio      NUMERIC,
    toFactor   NUMERIC,
    fromFactor NUMERIC
);

CREATE TABLE strategies (
    id    VARCHAR PRIMARY KEY,
    desc  TEXT,
    start TIMESTAMP DEFAULT timezone('utc', now())
);

CREATE TABLE positions (
    time    TIMESTAMP,
    stratid VARCHAR REFERENCES strategies(id),
    secid   VARCHAR REFERENCES stocks(id),
    shares  NUMERIC
);

CREATE INDEX ON positions(TIME DESC, strategy);

SELECT create_hypertable(
    'positions', 
    'time', 
    create_default_indexes => FALSE,
    chunk_time_interval => INTERVAL '1 day'
);
