package model

type Conversion struct {
	Amount       float64
	FromCurrency *Currency
	ToCurrency   *Currency
	Result       float64
}

func (c *Conversion) GetAmount() float64 {
	return c.Amount
}

func (c *Conversion) GetFromCurrency() *Currency {
	return c.FromCurrency
}

func (c *Conversion) GetToCurrency() *Currency {
	return c.ToCurrency
}

func (c *Conversion) GetResult() float64 {
	return c.Result
}
