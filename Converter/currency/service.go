package currency

import (
	"errors"
)

var (
	ErrInvalidAmount    = errors.New("amount должен быть положительным числом")
	ErrCurrencyNotfound = errors.New("такая валюта не найдена")
)

type Service struct {
	rates map[string]float64
}

func NewService() *Service {
	return &Service{
		rates: map[string]float64{
			"USD Доллар США":           1.0,
			"EUR Евро":                 0.92,
			"RUB Российский рубль":     90.0,
			"JPY Японская иена":        157.0,
			"CNY Китайский юань":       7.25,
			"GBP Британский фунт":      0.78,
			"KZT Казахстанский тенге":  460.0,
			"TRY Турецкая лира":        32.5,
			"INR Индийская рупия":      83.0,
			"BRL Бразильский реал":     5.12,
			"AUD Австралийский доллар": 1.50,
			"CAD Канадский доллар":     1.36,
			"CHF Швейцарский франк":    0.89,
			"SEK Шведская крона":       10.8,
			"NOK Норвежская крона":     10.5,
		},
	}
}

func (s *Service) Convert(req ConversionRequest) (ConversionResult, error) {
	if req.Amount <= 0 {
		return ConversionResult{}, ErrInvalidAmount
	}

	fromRate, ok := s.rates[req.From]
	if !ok {
		return ConversionResult{}, ErrCurrencyNotfound
	}

	toRate, ok := s.rates[req.To]
	if !ok {
		return ConversionResult{}, ErrCurrencyNotfound
	}

	rate := toRate / fromRate

	return ConversionResult{
		Amount:          req.Amount,
		ConvertedAmount: req.Amount * rate,
		From:            req.From,
		To:              req.To,
		Rate:            rate,
	}, nil
}

func (s *Service) ListCurrencies() []string {
	currencies := make([]string, 0, len(s.rates))
	for code := range s.rates {
		currencies = append(currencies, code)
	}
	return currencies
}
