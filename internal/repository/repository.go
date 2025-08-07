package repository

import (
	"fmt"
	"learnpack/src/currency-converter/internal/model"
	"log"
	"sync"
	"time"
)

var (
	currencyChan = make(chan *model.Currency)
	logChan      = make(chan logEntry)
	stopChan     = make(chan chan struct{})
)

var (
	currencies    []*model.Currency
	currencyMutex sync.Mutex
)

type logEntry struct {
	entityType string
	entities   []interface{}
}

func StoreEntity(entity model.Entity) {
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
			logChan <- logEntry{"Currency", []interface{}{currency}}
		case <-stopChan:
			return
		}
	}
}

func GetAllCurrencies() []*model.Currency {
	currencyMutex.Lock()
	defer currencyMutex.Unlock()
	return currencies
}

func startLogging() {
	var prevCurrencies = make(map[string]bool)

	for {
		time.Sleep(time.Millisecond * 200)

		currencyMutex.Lock()
		currentCurrencies := currencies
		currencyMutex.Unlock()

		for _, cur := range currentCurrencies {
			if !prevCurrencies[cur.Code] {
				log.Printf("Добавлена валюта: %v", cur)
				prevCurrencies[cur.Code] = true
			}
		}
	}
}

func InitRepository() {
	go processCurrencies()
	go startLogging()
}

func StopRepository() {
	close(stopChan)
	close(currencyChan)
	close(logChan)
}

func PrintRepositoryState() {
	currencyMutex.Lock()
	defer currencyMutex.Unlock()

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
