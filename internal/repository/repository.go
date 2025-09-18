package repository

import (
	"currency-converter/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	currencyFile   = "data/currency.json"
	conversionFile = "data/conversion.json"
)

type Repository interface {
	Store(entity model.Entity) error
	GetCurrencies() map[string]*model.Currency
	GetConversions() []*model.Conversion
	DeleteCurrency(code string) error
	UpdateCurrency(currency *model.Currency) error
	LoadCurrencies() error
	LoadConversions() error
}

type repo struct {
	mu          sync.RWMutex
	currencies  map[string]*model.Currency
	conversions []*model.Conversion
}

func NewRepository() Repository {
	return &repo{
		currencies:  make(map[string]*model.Currency),
		conversions: []*model.Conversion{},
	}
}

func (r *repo) Store(entity model.Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch v := entity.(type) {
	case *model.Currency:
		r.currencies[v.Code] = v
		return r.saveCurrenciesToFile()
	case *model.Conversion:
		r.conversions = append(r.conversions, v)
		return r.saveConversionsToFile()
	default:
		return fmt.Errorf("unknown entity type provided")
	}
}

func (r *repo) saveCurrenciesToFile() error {
	data, err := json.MarshalIndent(r.currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal currencies data: %w", err)
	}
	if err := os.WriteFile(currencyFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write currencies to file: %w", err)
	}
	return nil
}

func (r *repo) saveConversionsToFile() error {
	data, err := json.MarshalIndent(r.conversions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversions data: %w", err)
	}
	if err := os.WriteFile(conversionFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write conversions to file: %w", err)
	}
	return nil
}

func (r *repo) LoadCurrencies() error {
	fileData, err := os.ReadFile(currencyFile)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("data", 0755)
			return nil
		}
		return fmt.Errorf("failed to read currencies file: %w", err)
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if err := json.Unmarshal(fileData, &r.currencies); err != nil {
		return fmt.Errorf("failed to unmarshal currencies data: %w", err)
	}
	return nil
}

func (r *repo) LoadConversions() error {
	fileData, err := os.ReadFile(conversionFile)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll("data", 0755)
			return nil
		}
		return fmt.Errorf("failed to read conversions file: %w", err)
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if err := json.Unmarshal(fileData, &r.conversions); err != nil {
		return fmt.Errorf("failed to unmarshal conversions data: %w", err)
	}
	return nil
}

func (r *repo) DeleteCurrency(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.currencies[code]; !exists {
		return fmt.Errorf("currency %s not found in repository", code)
	}
	
	delete(r.currencies, code)
	return r.saveCurrenciesToFile()
}

func (r *repo) UpdateCurrency(currency *model.Currency) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.currencies[currency.Code]; !exists {
		return fmt.Errorf("currency %s not found for update", currency.Code)
	}
	
	r.currencies[currency.Code] = currency
	return r.saveCurrenciesToFile()
}

func (r *repo) GetCurrencies() map[string]*model.Currency {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copyMap := make(map[string]*model.Currency, len(r.currencies))
	for code, currency := range r.currencies {
		copyMap[code] = currency
	}
	return copyMap
}

func (r *repo) GetConversions() []*model.Conversion {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Conversion, len(r.conversions))
	copy(result, r.conversions)
	return result
}