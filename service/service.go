package service

import (
	"context"
	"fmt"
	"learnpack/src/currency-converter/internal/model"
	"learnpack/src/currency-converter/repository"
	"log"
	"time"
)

var (
	entityChan = make(chan model.Entity, 10)
)

// Слушает entityChan и отправляет в repository.Store
func processEntities(ctx context.Context, entityChan <-chan model.Entity) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nОбработчик получил сигнал об остановке.")
			return
		case e, ok := <-entityChan:
			if !ok {
				return
			}
			if e != nil {
				repository.Store(e)
			}
		}
	}
}

// Каждые 2 секунды создает фейк валюты и отправляет в entityChan
func generateData(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	i := 1
	for {
		select {
		case <-ticker.C:
			cur := model.NewCurrency(
				"CUR"+string(rune('A'+i)),
				float64(i)*1.1,
				"TestCurrency",
				"$",
			)
			entityChan <- cur
			i++
		case <-ctx.Done():
			fmt.Println("\nГенератор получил сигнал\nоб остановке и закрывает канал")
			close(entityChan)
			return
		}
	}
}

// Каждые 200мс проверяем текущую копию мапы на наличие новой валюты
// и печатаем в log
func startLogging(ctx context.Context) {
	seen := make(map[string]bool)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nЛоггер получил сигнал об остановке.")
			return
		case <-time.After(200 * time.Millisecond):
			current := repository.GetCurrencies()
			for _, cur := range current {
				if !seen[cur.Code()] {
					log.Printf("Добавлена валюта: %s (%s)", cur.Code(), cur.Name())
					seen[cur.Code()] = true //Пометка, как выделенная
				}
			}
		}
	}
}

func InitService(ctx context.Context) {
	go processEntities(ctx, entityChan)
	go generateData(ctx)
	go startLogging(ctx)
}
