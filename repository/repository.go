package repository

import (
	"currency-converter/internal/model"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	currencyFile   = "data/currency.json"
	conversionFile = "data/conversion.json"
)

var (
	mu sync.RWMutex

	currencies  = make(map[string]*model.Currency) //Для O(1)
	conversions []*model.Conversion                //Для истории конверсий
)

// Сравнивает сущность по типу, полученную из processEntities
// и отправляет в мапу или слайс и свои json файлы.
func Store(e model.Entity) {
	mu.Lock()
	defer mu.Unlock()

	switch v := e.(type) {
	case *model.Currency:
		currencies[v.Code] = v //(код валюты как ключ)
		if err := SaveCurToFile(); err != nil {
			log.Println("error saving currencies:", err)
		}
	case *model.Conversion:
		conversions = append(conversions, v)
		if err := SaveConvToFile(); err != nil {
			log.Println("error saving conversions:", err)
		}
	}
}

// Содержимое мапы в json
func SaveCurToFile() error {
	data, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling currencies: %w", err)
	}
	if err := os.WriteFile(currencyFile, data, 0644); err != nil {
		return fmt.Errorf("error writing the currencys file: %w", err)
	}
	return nil
}

// Содердимове слайса в json
func SaveConvToFile() error {
	data, err := json.MarshalIndent(conversions, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling conversions: %w", err)
	}
	if err := os.WriteFile(conversionFile, data, 0644); err != nil {
		return fmt.Errorf("error writing the conversions file: %w", err)
	}
	return nil
}

// Загрузка содержимого json в мапу
func LoadCurrenciesFromFile() error {
	fileData, err := os.ReadFile(currencyFile)
	if err != nil {
		return fmt.Errorf("error reading the currencys file: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if err := json.Unmarshal(fileData, &currencies); err != nil {
		return fmt.Errorf("error unmarshaling the currencys file: %w", err)
	}

	return nil
}

// Загрузка содержимого json в слайс
func LoadConversionsFromFile() error {
	fileData, err := os.ReadFile(conversionFile)
	if err != nil {
		return fmt.Errorf("error reading the conversions file: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if err := json.Unmarshal(fileData, &conversions); err != nil {
		//При первом запуске будет ошибка, т.к. json пуст
		return fmt.Errorf("error unmarshaling the conversions file: %w", err)
	}

	return nil
}

// Функция для удаления валюты(не копии)
func DeleteCurFromMap(code string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := currencies[code]; !ok {
		return fmt.Errorf("the currency %s not found", code)
	}
	delete(currencies, code)
	return SaveCurToFile()
}

// Функция для обновления валюты(не копии)
func UpdateCurInMap(cur *model.Currency) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := currencies[cur.Code]; !ok {
		return fmt.Errorf("the currency %s not found", cur.Code)
	}
	currencies[cur.Code] = cur
	return SaveCurToFile()
}

// Возвращает копию мапы валют
func GetCurrencies() map[string]*model.Currency {
	mu.RLock()
	defer mu.RUnlock()

	copyMap := make(map[string]*model.Currency, len(currencies))
	for k, v := range currencies {
		copyMap[k] = v
	}
	return copyMap
}

// Возвращает копию слайса конвертаций
func GetConversions() []*model.Conversion {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]*model.Conversion, len(conversions))
	copy(result, conversions)
	return result
}
