package model

type Conversion struct {
	amount float64
	from   *Currency
	to     *Currency
	result float64
}

// Конструктор конвертирования
func NewConversion(amount float64, from *Currency, to *Currency, result float64) *Conversion {
	return &Conversion{
		amount: amount,
		from:   from,
		to:     to,
		result: result,
	}
}

//Геттеры для приватных полей структуры 
func (c *Conversion) Amount() float64 {
	return c.amount
}

func (c *Conversion) FromCurrency() *Currency {
	return c.from
}

func (c *Conversion) ToCurrency() *Currency {
	return c.to
}

func (c *Conversion) Result() float64 {
	return c.result
}
