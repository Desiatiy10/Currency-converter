package repository

import (
	"learnpack/src/currency-converter/internal/model"
)

type LogEntry struct {
	EntityType string
	Entities   []interface{}
}

func ProcessEntities(storeFunc func(model.Entity)) {
	currencyUSD := model.Currency{
		Code:   "USD",
		Rate:   1.0,
		Name:   "US Dollar",
		Symbol: "$",
	}

	storeFunc(&currencyUSD)
}
