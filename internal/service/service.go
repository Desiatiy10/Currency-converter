package service

import (
	"learnpack/src/currency-converter/internal/model"
	"learnpack/src/currency-converter/internal/repository"
	"time"
)

func InitService() {

	repository.InitRepository()

	go ProcessEntities()

	time.Sleep(1 * time.Second)

	repository.PrintRepositoryState()
}

func ProcessEntities() {
	currencyUSD := model.Currency{
		Code:   "USD",
		Rate:   1.0,
		Name:   "US Dollar",
		Symbol: "$",
	}

	currencyEUR := model.Currency{
		Code:   "EUR",
		Rate:   0.85,
		Name:   "Euro",
		Symbol: "€",
	}

	repository.StoreEntity(&currencyUSD)
	repository.StoreEntity(&currencyEUR)

}

func StopService() {
	repository.StopRepository()
}
