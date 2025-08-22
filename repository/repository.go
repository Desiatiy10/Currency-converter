package repository

import (
	"learnpack/src/currency-converter/internal/model"
	"sync"
)

var (
	mu sync.RWMutex

	currencies  = make(map[string]*model.Currency) //Для O(1)
	conversions []*model.Conversion                //Для истории конверсий
)

//Сравнивает сущность по типу, полученную из processEntities
//и отправляет в мапу для валюты или слайс для конвертирования.
func Store(e model.Entity) {
	mu.Lock()
	defer mu.Unlock()

	switch v := e.(type) {

	case *model.Currency:
		currencies[v.Code()] = v //(код валюты как ключ)

	case *model.Conversion:
		conversions = append(conversions, v)
	}
}

//Возвращает копию мапы валют
func GetCurrencies() map[string]*model.Currency {
	mu.RLock()
	defer mu.RUnlock()

	copyMap := make(map[string]*model.Currency, len(currencies))
	for k, v := range currencies {
		copyMap[k] = v
	}

	return copyMap
}

//Возвращает копию слайса конвертаций
func GetConversions() []*model.Conversion {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]*model.Conversion, len(conversions))
	copy(result, conversions)
	return result
}
