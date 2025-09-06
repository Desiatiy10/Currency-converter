package model

type Conversion struct {
	Amount float64   `json:"amount"`
	From   *Currency `json:"from"`
	To     *Currency `json:"to"`
	Result float64   `json:"result"`
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
