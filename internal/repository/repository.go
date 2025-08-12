package repository

import (
	"learnpack/src/currency-converter/internal/model"
	"sync"
)

var (
	Currencies    = make(map[string]*model.Currency)
	CurrencyMutex sync.RWMutex
)

type LogEntry struct {
	EntityType string
	Entities   []interface{}
}

func GetAllCurrencies()map[string]*model.Currency {
	CurrencyMutex.RLock()
	defer CurrencyMutex.RUnlock()

	copyMap := make(map[string]*model.Currency, len(Currencies))
	for k, v := range Currencies {
		copyMap[k] = v
	}
	return copyMap
}

func AddCurrency(currency *model.Currency) {
	CurrencyMutex.Lock()
	defer CurrencyMutex.Unlock()

	Currencies[currency.Code] = currency
}

func ProcessEntities(storeFunc func(model.Entity)) { 
}
