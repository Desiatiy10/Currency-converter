package model

type Currency struct {
	Code   string
	Rate   float64
	Name   string
	Symbol string
}

func (c *Currency) GetCode() string {
	return c.Code
}

func (c *Currency) GetRate() float64 {
	return c.Rate
}

func (c *Currency) GetName() string {
	return c.Name
}

func (c *Currency) GetSymbol() string {
	return c.Symbol
}
