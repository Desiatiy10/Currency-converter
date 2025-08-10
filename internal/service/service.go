package service

import (
	"fmt"
	mod "learnpack/src/currency-converter/internal/model"
	repo "learnpack/src/currency-converter/internal/repository"
	"log"
	"time"
)

var (
	currencyChan = make(chan *mod.Currency)
	logChan      = make(chan repo.LogEntry)
	stopChan     = make(chan struct{})
)

func storeEntity(entity mod.Entity) {
	switch v := entity.(type) {
	case *mod.Currency:
		currencyChan <- v
	default:
		log.Panicf("неизвестный тип: %T", v)
	}
}

func processCurrencies() {
	select {
	case currency := <-currencyChan:
		repo.AddCurrency(currency)
		logChan <- repo.LogEntry{
			EntityType: "Currency",
			Entities:   []interface{}{currency},
		}
	case <-stopChan:
		return
	}
}

func startLogging() {
	var prevCurrencies = make(map[string]bool)
	for {
		select {
		case <-stopChan:
			return
		default:
			time.Sleep(time.Millisecond * 200)
		}

		currentCurrencies := repo.GetAllCurrencies()

		for _, cur := range currentCurrencies {
			if !prevCurrencies[cur.Code] {
				log.Printf("Добавлена валюта: %v", cur)
				prevCurrencies[cur.Code] = true
			}
		}
	}
}

func InitRepository() {
	go func() {
		repo.GetAllCurrencies()
		repo.ProcessEntities(storeEntity)
	}()
	go startLogging()
}

func InitService() {
	go processCurrencies()
	InitRepository()

	time.Sleep(time.Second * 1)

	printRepositoryState()
}

func StopRepository() {
	close(stopChan)
	close(currencyChan)
	close(logChan)
}

func StopService() {
	StopRepository()
}

func printRepositoryState() {
	currencies := repo.GetAllCurrencies()
	if currencies == nil {
		log.Println("Ошибка: список валют не получен")
		return
	}

	fmt.Println("\nТекущее состояние хранилища:")

	fmt.Println("\nВалюты в системе:")
	for i, currency := range currencies {
		fmt.Printf("Валюта №%d:\n", i+1)
		fmt.Printf("  Код: %s\n", currency.Code)
		fmt.Printf("  Название: %s\n", currency.Name)
		fmt.Printf("  Символ: %s\n", currency.Symbol)
		fmt.Printf("  Курс: %f\n", currency.Rate)
	}

	fmt.Printf("\nВсего валют: %d\n", len(currencies))
}
