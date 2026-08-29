package repository

import (
	"currency-converter/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var dataDir = findDataDir()

type Repository interface {
	Store(entity model.Entity) error
	GetCurrencies() map[string]*model.Currency
	GetConversions() []*model.Conversion
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

func findDataDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return "data"
	}

	for {
		if _, err := os.Stat(filepath.Join(workingDir, "go.mod")); err == nil {
			return filepath.Join(workingDir, "data")
		}

		parentDir := filepath.Dir(workingDir)
		if parentDir == workingDir {
			return filepath.Join(workingDir, "data")
		}
		workingDir = parentDir
	}
}

func currencyFile() string {
	return filepath.Join(dataDir, "currency.json")
}

func conversionFile() string {
	return filepath.Join(dataDir, "conversion.json")
}

func ensureDataDir() error {
	return os.MkdirAll(dataDir, 0755)
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
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	data, err := json.MarshalIndent(r.currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal currencies data: %w", err)
	}
	if err := os.WriteFile(currencyFile(), data, 0644); err != nil {
		return fmt.Errorf("failed to write currencies to file: %w", err)
	}
	return nil
}

func (r *repo) saveConversionsToFile() error {
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	data, err := json.MarshalIndent(r.conversions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversions data: %w", err)
	}
	if err := os.WriteFile(conversionFile(), data, 0644); err != nil {
		return fmt.Errorf("failed to write conversions to file: %w", err)
	}
	return nil
}

func (r *repo) LoadCurrencies() error {
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to prepare data directory: %w", err)
	}

	fileData, err := os.ReadFile(currencyFile())
	if err != nil {
		if os.IsNotExist(err) {
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
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to prepare data directory: %w", err)
	}

	fileData, err := os.ReadFile(conversionFile())
	if err != nil {
		if os.IsNotExist(err) {
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
