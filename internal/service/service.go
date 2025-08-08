package service

import (
	"fmt"
	"learnpack/src/currency-converter/internal/model"
	"learnpack/src/currency-converter/internal/repository"
	"log"
	"sync"
	"time"
)

var (
	currencyChan = make(chan *model.Currency)
	logChan      = make(chan repository.LogEntry)
	stopChan     = make(chan struct{})

	currencies    []*model.Currency
	currencyMutex sync.Mutex
)

func getAllCurrencies() []*model.Currency {
	currencyMutex.Lock()
	defer currencyMutex.Unlock()
	return currencies
}

func storeEntity(entity model.Entity) {
	switch v := entity.(type) {
	case *model.Currency:
		currencyChan <- v
	default:
		log.Panicf("неизвестный тип: %T", v)
	}
}

func processCurrencies() {
	for {
		select {
		case currency := <-currencyChan:
			currencyMutex.Lock()
			currencies = append(currencies, currency)
			currencyMutex.Unlock()
			logChan <- repository.LogEntry{
				EntityType: "Currency",
				Entities:   []interface{}{currency}}
		case <-stopChan:
			return
		}
	}
}

func startLogging() {
	var (
		prevCurrencies = make(map[string]bool)
	)
	for {
		time.Sleep(time.Millisecond * 200)

		currentCurrencies := getAllCurrencies()

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
		getAllCurrencies()
		repository.ProcessEntities(storeEntity)
	}()
	go startLogging()
}

func InitService() {
	go processCurrencies()
	InitRepository()

	time.Sleep(time.Second)

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
	currencies := getAllCurrencies()
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
