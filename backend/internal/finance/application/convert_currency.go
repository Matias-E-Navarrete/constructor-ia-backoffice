package application

import (
	"context"
	"errors"
)

var ErrUnknownCurrency = errors.New("finance: unknown currency")

type ConvertCurrency struct {
	Rates *ExchangeRateProvider
}

type ConversionResult struct {
	ConvertedAmount float64
	Rate            float64
}

// Execute converts `amount` of `from` currency into `to` currency using the
// current live (or last-known-good) USD-based rate table.
func (uc *ConvertCurrency) Execute(ctx context.Context, amount float64, from, to string) (ConversionResult, error) {
	rates := uc.Rates.Rates(ctx)

	fromRate, ok := rates[from]
	if !ok {
		return ConversionResult{}, ErrUnknownCurrency
	}
	toRate, ok := rates[to]
	if !ok {
		return ConversionResult{}, ErrUnknownCurrency
	}

	usd := amount / fromRate
	converted := usd * toRate
	rate := toRate / fromRate
	return ConversionResult{ConvertedAmount: converted, Rate: rate}, nil
}
