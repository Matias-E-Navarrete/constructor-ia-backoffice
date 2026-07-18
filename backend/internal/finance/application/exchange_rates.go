package application

import (
	"context"
	"errors"
)

var ErrUnknownCurrency = errors.New("finance: unknown currency")

// demoRatesToUSD are illustrative, fixed exchange rates (units of currency
// per 1 USD), standing in for a real live-rates provider. Good enough to
// demonstrate "convert as you type" without an external API dependency.
var demoRatesToUSD = map[string]float64{
	"USD": 1,
	"EUR": 0.92,
	"GBP": 0.79,
	"CLP": 950,
	"BRL": 5.4,
	"MXN": 17.0,
	"ARS": 1000,
}

type ConvertCurrency struct{}

type ConversionResult struct {
	ConvertedAmount float64
	Rate            float64
}

// Execute converts `amount` of `from` currency into `to` currency.
func (uc *ConvertCurrency) Execute(ctx context.Context, amount float64, from, to string) (ConversionResult, error) {
	fromRate, ok := demoRatesToUSD[from]
	if !ok {
		return ConversionResult{}, ErrUnknownCurrency
	}
	toRate, ok := demoRatesToUSD[to]
	if !ok {
		return ConversionResult{}, ErrUnknownCurrency
	}
	usd := amount / fromRate
	converted := usd * toRate
	rate := toRate / fromRate
	return ConversionResult{ConvertedAmount: converted, Rate: rate}, nil
}
