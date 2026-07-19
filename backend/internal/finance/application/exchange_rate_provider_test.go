package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExchangeRateProvider_FetchesAndCachesLiveRates(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":"success","rates":{"USD":1,"EUR":0.9,"CLP":900}}`))
	}))
	defer server.Close()

	p := &ExchangeRateProvider{BaseURL: server.URL, Client: server.Client(), CacheTTL: time.Hour}

	rates := p.Rates(context.Background())
	if rates["EUR"] != 0.9 || rates["CLP"] != 900 {
		t.Fatalf("unexpected rates from live fetch: %+v", rates)
	}
	if calls != 1 {
		t.Fatalf("expected 1 upstream call, got %d", calls)
	}

	// A second call within the TTL should be served from cache, not hit the
	// upstream server again.
	p.Rates(context.Background())
	if calls != 1 {
		t.Fatalf("expected cache to avoid a second upstream call, got %d calls", calls)
	}
}

func TestExchangeRateProvider_FallsBackWhenUnreachable(t *testing.T) {
	p := &ExchangeRateProvider{BaseURL: "http://127.0.0.1:1", Client: &http.Client{Timeout: time.Second}, CacheTTL: time.Hour}

	rates := p.Rates(context.Background())
	if rates["USD"] != 1 {
		t.Fatalf("expected fallback table to be used, got %+v", rates)
	}
}

func TestExchangeRateProvider_FallsBackToStaleCacheOnFailure(t *testing.T) {
	up := true
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !up {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":"success","rates":{"USD":1,"EUR":0.87}}`))
	}))
	defer server.Close()

	p := &ExchangeRateProvider{BaseURL: server.URL, Client: server.Client(), CacheTTL: time.Millisecond}

	first := p.Rates(context.Background())
	if first["EUR"] != 0.87 {
		t.Fatalf("expected live rate on first call, got %+v", first)
	}

	time.Sleep(2 * time.Millisecond) // let the cache go stale
	up = false

	second := p.Rates(context.Background())
	if second["EUR"] != 0.87 {
		t.Fatalf("expected stale cache to be reused when upstream fails, got %+v", second)
	}
}

func TestConvertCurrency_UsesProviderRates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":"success","rates":{"USD":1,"EUR":0.5,"GBP":0.25}}`))
	}))
	defer server.Close()

	uc := &ConvertCurrency{Rates: &ExchangeRateProvider{BaseURL: server.URL, Client: server.Client(), CacheTTL: time.Hour}}

	result, err := uc.Execute(context.Background(), 10, "EUR", "GBP")
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	// 10 EUR -> 20 USD -> 5 GBP
	if result.ConvertedAmount != 5 {
		t.Fatalf("expected 5, got %f", result.ConvertedAmount)
	}

	_, err = uc.Execute(context.Background(), 10, "XXX", "GBP")
	if err != ErrUnknownCurrency {
		t.Fatalf("expected ErrUnknownCurrency, got %v", err)
	}
}
