package models

type ExchangeRatesResponse struct {
	Base          string                   `json:"base"`
	Rates         map[string]float64       `json:"rates,omitempty"`
	AnalyzedRates map[string]AnalyzedRates `json:"rates_analyze,omitempty"`
	Date          string                   `json:"date,omitempty"`
}

type AnalyzedRates struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Count int     `json:"-"`
	Sum   float64 `json:"-"`
}
