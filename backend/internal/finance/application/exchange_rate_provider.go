package application

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

// fallbackRatesToUSD are illustrative, fixed rates (units of currency per 1
// USD) used only when the live provider can't be reached, so conversion
// degrades gracefully instead of failing outright.
var fallbackRatesToUSD = map[string]float64{
	"USD": 1,
	"EUR": 0.92,
	"GBP": 0.79,
	"CLP": 950,
	"BRL": 5.4,
	"MXN": 17.0,
	"ARS": 1000,
}

const defaultRatesCacheTTL = time.Hour

type ratesProviderResponse struct {
	Result string             `json:"result"`
	Rates  map[string]float64 `json:"rates"`
}

// ExchangeRateProvider fetches live USD-based exchange rates from a free,
// keyless provider and caches them in memory for CacheTTL, so a busy wallet
// page doesn't hammer the upstream API on every keystroke. It never returns
// an error to callers: an unreachable provider falls back to the last good
// cache, or to a static illustrative table as a last resort.
type ExchangeRateProvider struct {
	BaseURL  string
	Client   *http.Client
	CacheTTL time.Duration

	mu        sync.Mutex
	cached    map[string]float64
	fetchedAt time.Time
}

func NewExchangeRateProvider() *ExchangeRateProvider {
	return &ExchangeRateProvider{
		BaseURL:  "https://open.er-api.com/v6/latest/USD",
		Client:   &http.Client{Timeout: 5 * time.Second},
		CacheTTL: defaultRatesCacheTTL,
	}
}

func (p *ExchangeRateProvider) Rates(ctx context.Context) map[string]float64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cached != nil && time.Since(p.fetchedAt) < p.CacheTTL {
		return p.cached
	}

	fresh, err := p.fetch(ctx)
	if err != nil {
		if p.cached != nil {
			return p.cached
		}
		return fallbackRatesToUSD
	}

	p.cached = fresh
	p.fetchedAt = time.Now()
	return p.cached
}

func (p *ExchangeRateProvider) fetch(ctx context.Context) (map[string]float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("finance: rates provider returned a non-200 status")
	}

	var parsed ratesProviderResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if parsed.Result != "success" || len(parsed.Rates) == 0 {
		return nil, errors.New("finance: rates provider returned no usable rates")
	}
	return parsed.Rates, nil
}
