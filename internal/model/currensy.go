package model

type Currency struct {
	code   string
	rate   float64
	name   string
	symbol string
}

// Конструктор новой валюты
func NewCurrency(code string, rate float64, name string, symbol string) *Currency {
	return &Currency{
		code:   code,
		rate:   rate,
		name:   name,
		symbol: symbol,
	}
}

// Геттеры для приватных полей структуры
func (c *Currency) Code() string {
	return c.code
}

func (c *Currency) Rate() float64 {
	return c.rate
}

func (c *Currency) Name() string {
	return c.name
}

func (c *Currency) Symbol() string {
	return c.symbol
}
