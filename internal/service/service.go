package service

import (
	"learnpack/src/currency-converter/internal/model"
	repo "learnpack/src/currency-converter/internal/repository"
	"log"
	"time"
)

var (
	currencyChan = make(chan *model.Currency)
	logChan      = make(chan repo.LogEntry)
	stopChan     = make(chan struct{})
)

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
			repo.AddCurrency(currency)
			select {
			case logChan <- repo.LogEntry{
				EntityType: "Currency",
				Entities:   []interface{}{currency},
			}:
			case <-stopChan:
				return
			}
		case <-stopChan:
			return
		}
	}
}

func startLogging() {
	prevCurrencies := make(map[string]bool)

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
		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				repo.ProcessEntities(storeEntity)
			case <-stopChan:
				return
			}
		}
	}()
	go startLogging()
}

func InitService() {
	go processCurrencies()
	InitRepository()
}

func StopRepository() {
	close(stopChan)
	close(currencyChan)
	close(logChan)
}

func StopService() {
	StopRepository()
}
