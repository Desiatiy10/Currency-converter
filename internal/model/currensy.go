package model

type Currency struct {
	Code   string  `json:"code"`
	Rate   float64 `json:"rate"`
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
}

// Конструктор новой валюты
func NewCurrency(code string, rate float64, name string, symbol string) *Currency {
	return &Currency{
		Code:   code,
		Rate:   rate,
		Name:   name,
		Symbol: symbol,
	}
}


