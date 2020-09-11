package iex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"golang.org/x/time/rate"
)

// APIConnection generalizes a http client
type APIConnection struct {
	rateLimiter *rate.Limiter
	baseURL     url.URL
	apiKey      string
}

// NewAPIConnection creates a http client with personal api key
func NewAPIConnection(key string) *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(Per(100, time.Second), 1),
		baseURL: url.URL{
			Scheme: "https",
			Host:   "sandbox.iexapis.com",
			Path:   "stable/",
		},
		apiKey: key,
	}
}

// Per tracks events per unit of time
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

// AllStocks returns all active stocks and their accompanied data
func (a *APIConnection) AllStocks(ctx context.Context) ([]Stock, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	endpoint := a.baseURL.ResolveReference(&url.URL{Path: "ref-data/symbols"})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = qparams.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stks []Stock
	if err := json.NewDecoder(resp.Body).Decode(&stks); err != nil {
		return nil, err
	}

	// My Custom Filter: Possibly make this a functional option
	re := regexp.MustCompile(`cs|ad`)
	var fstks []Stock
	for _, s := range stks {
		if len(s.Cik) == 0 || !re.MatchString(s.Type) {
			continue
		}
		fstks = append(fstks, s)
	}
	return fstks, nil
}

// Prices returns the historical prices for a stock
func (a *APIConnection) Prices(ctx context.Context, symbol, lookback string) (*PriceHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "chart", lookback)
	endpoint := a.baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = qparams.Encode()
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prices []Prices
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}
	priceHist := PriceHistory{
		Symbol: symbol,
		Prices: prices,
	}
	return &priceHist, nil
}

// Dividends returns the historical dividend information for a stock
func (a *APIConnection) Dividends(ctx context.Context, symbol, lookback string) (*DividendHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "dividends", lookback)
	endpoint := a.baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = qparams.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dividends []Dividend
	if err := json.NewDecoder(resp.Body).Decode(&dividends); err != nil {
		return nil, err
	}
	divHist := DividendHistory{
		Symbol:    symbol,
		Dividends: dividends,
	}
	return &divHist, nil
}

// Splits returns the historical split information for a stock
func (a *APIConnection) Splits(ctx context.Context, symbol, lookback string) (*SplitHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "splits", lookback)
	endpoint := a.baseURL.ResolveReference(&url.URL{Path: urlpath})
	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	request.URL.RawQuery = qparams.Encode()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var splits []Split
	if err := json.NewDecoder(resp.Body).Decode(&splits); err != nil {
		return nil, err
	}
	splitHist := SplitHistory{
		Symbol: symbol,
		Splits: splits,
	}
	return &splitHist, nil
}
