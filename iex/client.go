package iex

import (
	"context"
	"encoding/json"
	"log"
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
	lookback    string
}

// NewAPIConnection creates a http client with personal api key
func NewAPIConnection(host, key, lookback string, duration time.Duration) *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(Per(1, duration), 1),
		baseURL: url.URL{
			Scheme: "https",
			Host:   host,
			Path:   "stable/",
		},
		apiKey:   key,
		lookback: lookback,
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
func (a *APIConnection) Prices(ctx context.Context, symbol string) (*PriceHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "chart", a.lookback)
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
	log.Printf("price status: %s\n", resp.Status)

	var ph PriceHistory
	ph.Symbol = symbol
	if err := json.NewDecoder(resp.Body).Decode(&ph.Prices); err != nil {
		return nil, err
	}
	return &ph, nil
}

// Dividends returns the historical dividend information for a stock
func (a *APIConnection) Dividends(ctx context.Context, symbol string) (*DividendHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "dividends", a.lookback)
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
	log.Printf("dividend status: %s\n", resp.Status)

	var dh DividendHistory
	dh.Symbol = symbol
	if err := json.NewDecoder(resp.Body).Decode(&dh.Dividends); err != nil {
		return nil, err
	}
	return &dh, nil
}

// Splits returns the historical split information for a stock
func (a *APIConnection) Splits(ctx context.Context, symbol string) (*SplitHistory, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	qparams := make(url.Values)
	qparams.Set("token", a.apiKey)
	urlpath := path.Join("stock", symbol, "splits", a.lookback)
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

	log.Printf("split status: %s\n", resp.Status)
	var sh SplitHistory
	sh.Symbol = symbol
	if err := json.NewDecoder(resp.Body).Decode(&sh.Splits); err != nil {
		return nil, err
	}
	return &sh, nil
}
