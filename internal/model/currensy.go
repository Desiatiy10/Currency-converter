package model

type Currency struct {
	Code   string  `json:"Код"`
	Rate   float64 `json:"Курс"`
	Name   string  `json:"Название"`
	Symbol string  `json:"Символ"`
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
