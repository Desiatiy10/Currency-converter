package repository

import (
	"learnpack/src/currency-converter/internal/model"
	"sync"
)

var (
	Currencies    []*model.Currency
	CurrencyMutex sync.Mutex
)

type LogEntry struct {
	EntityType string
	Entities   []interface{}
}

func GetAllCurrencies() []*model.Currency {
	CurrencyMutex.Lock()
	defer CurrencyMutex.Unlock()
	return Currencies
}

func AddCurrency(currency *model.Currency) {
	CurrencyMutex.Lock()
	defer CurrencyMutex.Unlock()

	Currencies = append(Currencies, currency)
}

func ProcessEntities(storeFunc func(model.Entity)) {

}
