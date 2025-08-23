package model

type Conversion struct {
	Amount float64   `json:"Сумма"`
	From   *Currency `json:"Из "`
	To     *Currency `json:"В "`
	Result float64   `json:"Результат: "`
}

// Конструктор конвертирования
func NewConversion(amount float64, from *Currency, to *Currency, result float64) *Conversion {
	return &Conversion{
		Amount: amount,
		From:   from,
		To:     to,
		Result: result,
	}
}
