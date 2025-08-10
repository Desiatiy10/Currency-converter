package repository

import (
	mod "learnpack/src/currency-converter/internal/model"
	"log"
	"sync"
)

type LogEntry struct {
	EntityType string
	Entities   []interface{}
}

func GetAllCurrencies() []*mod.Currency {
	CurrencyMutex.Lock()
	defer CurrencyMutex.Unlock()
	return Currencies
}

func AddCurrency(currency *mod.Currency) {
    CurrencyMutex.Lock()
    defer CurrencyMutex.Unlock()

	for _, c := range Currencies {
		if c.Code == currency.Code {
			log.Printf("Валюта с кодом %s уже существует", currency.Code)
			return
		}
	}

    Currencies = append(Currencies, currency)
}

var (
	Currencies    []*mod.Currency
	CurrencyMutex sync.Mutex
)

func ProcessEntities(storeFunc func(mod.Entity)) {
	currencyUSD := mod.Currency{
		Code:   "USD",
		Rate:   1.0,
		Name:   "US Dollar",
		Symbol: "$",
	}

	storeFunc(&currencyUSD)
}
