package iex

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "cloud.iexapis.com",
	Path:   "stable/",
}

// Client generalizes a http client
type Client struct {
	c      *http.Client
	apiKey string
}

// NewClient creates a http client with personal api key
func NewClient() *Client {
	return &Client{
		c:      &http.Client{Timeout: time.Minute},
		apiKey: os.Getenv("IEXCLOUD_SECRET"),
	}
}

// AllStocks returns all active stocks and their accompanied data
func (c *Client) AllStocks() ([]Stock, error) {
	queryparams := make(url.Values)
	queryparams.Set("token", c.apiKey)
	endpoint := baseURL.ResolveReference(&url.URL{Path: "ref-data/symbols"})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = queryparams.Encode()
	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var stks []Stock
	if err := json.NewDecoder(resp.Body).Decode(&stks); err != nil {
		return nil, err
	}
	return stks, nil
}

// Prices returns the historical prices for a stock
func (c *Client) Prices(symbol, lookback string) (PriceHistory, error) {
	queryparams := make(url.Values)
	queryparams.Set("token", c.apiKey)
	urlpath := path.Join("stock", symbol, "chart", lookback)
	endpoint := baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = queryparams.Encode()
	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var prices PriceHistory
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}
	return prices, nil
}

// Dividends returns the historical dividend information for a stock
func (c *Client) Dividends(symbol, lookback string) (Dividend, error) {
	queryparams := make(url.Values)
	queryparams.Set("token", c.apiKey)
	urlpath := path.Join("stock", symbol, "dividends", lookback)
	endpoint := baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = queryparams.Encode()
	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var dividends Dividend
	if err := json.NewDecoder(resp.Body).Decode(&dividends); err != nil {
		return nil, err
	}
	return dividends, nil
}

// Splits returns the historical split information for a stock
func (c *Client) Splits(symbol, lookback string) (Split, error) {
	queryparams := make(url.Values)
	queryparams.Set("token", c.apiKey)
	urlpath := path.Join("stock", symbol, "splits", lookback)
	endpoint := baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = queryparams.Encode()
	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var splits Split
	if err := json.NewDecoder(resp.Body).Decode(&splits); err != nil {
		return nil, err
	}
	return splits, nil
}
