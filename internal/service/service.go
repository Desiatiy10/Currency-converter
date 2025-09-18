package service

import (
	"context"
	"currency-converter/internal/model"
	"currency-converter/internal/repository"
	"fmt"
	"log"
	"time"
)

type Service interface {
	AddEntity(e model.Entity) error

	CreateCurrency(*model.Currency) (*model.Currency, error)
	ListCurrencies() (map[string]*model.Currency, error)
	GetCurrency(code string) (*model.Currency, error)
	UpdateCurrency(cur *model.Currency) (*model.Currency, error)
	DeleteCurrency(code string) error

	ListConversions() ([]*model.Conversion, error)
	CreateConversion(amount float64, fromCode, toCode string) (*model.Conversion, error)
}

type service struct {
	repo       repository.Repository
	entityChan chan model.Entity
}

func NewService(repo repository.Repository) *service {
	return &service{
		repo:       repo,
		entityChan: make(chan model.Entity, 10),
	}
}

func (s *service) processEntities(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Entity processing stopped: context cancelled")
			return
		case entity, ok := <-s.entityChan:
			if !ok {
				log.Println("Entity channel closed")
				return
			}
			if entity != nil {
				err := s.repo.Store(entity)
				if err != nil {
					log.Printf("Failed to store entity: %v", err)
				}
			}
		}
	}
}

func (s *service) generateData(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	i := 1
	for {
		select {
		case <-ticker.C:
			currency := model.NewCurrency(
				"CUR"+string(rune('A'+i)),
				float64(i)*1.1,
				"TestCurrency",
				"$",
			)
			s.AddEntity(currency)
			log.Printf("Generated test currency: %s", currency.Code)
			i++
		case <-ctx.Done():
			close(s.entityChan)
			log.Println("Data generation stopped: context cancelled")
			return
		}
	}
}

func (s *service) startLogging(ctx context.Context) {
	seen := make(map[string]bool)
	for code := range s.repo.GetCurrencies() {
		seen[code] = true
	}
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Currency monitoring stopped: context cancelled")
			return
		case <-time.After(200 * time.Millisecond):
			currenciesData := s.repo.GetCurrencies()
			for _, cur := range currenciesData {
				if !seen[cur.Code] {
					log.Printf("New currency detected: %s - %s (Rate: %.2f)", cur.Code, cur.Name, cur.Rate)
					seen[cur.Code] = true
				}
			}
		}
	}
}

func InitService(ctx context.Context, repo repository.Repository) *service {
	s := NewService(repo)

	go s.processEntities(ctx)
	go s.generateData(ctx)
	go s.startLogging(ctx)

	log.Println("Currency converter service initialized successfully")
	return s
}

func (s *service) AddEntity(entity model.Entity) error {
	if entity == nil {
		return fmt.Errorf("cannot add nil entity")
	}
	select {
	case s.entityChan <- entity:
		return nil
	default:
		return fmt.Errorf("entity channel is full - cannot process request")
	}
}

func (s *service) CreateCurrency(cur *model.Currency) (*model.Currency, error) {
	if cur.Code == "" || cur.Rate <= 0 || cur.Name == "" || cur.Symbol == "" {
		return nil, fmt.Errorf("invalid currency data: all fields must be provided and rate must be positive")
	}
	
	if err := s.repo.Store(cur); err != nil {
		return nil, fmt.Errorf("failed to create currency: %v", err)
	}
	
	log.Printf("Currency created successfully: %s (%s)", cur.Code, cur.Name)
	return cur, nil
}

func (s *service) ListCurrencies() (map[string]*model.Currency, error) {
	currencies := s.repo.GetCurrencies()
	log.Printf("Retrieved %d currencies from repository", len(currencies))
	return currencies, nil
}

func (s *service) GetCurrency(code string) (*model.Currency, error) {
	if code == "" {
		return nil, fmt.Errorf("currency code cannot be empty")
	}
	
	data := s.repo.GetCurrencies()
	if cur, ok := data[code]; ok {
		log.Printf("Currency found: %s", code)
		return cur, nil
	}
	
	return nil, fmt.Errorf("currency '%s' not found in the system", code)
}

func (s *service) UpdateCurrency(cur *model.Currency) (*model.Currency, error) {
	if cur.Code == "" {
		return nil, fmt.Errorf("currency code is required for update")
	}
	
	err := s.repo.UpdateCurrency(cur)
	if err != nil {
		return nil, fmt.Errorf("failed to update currency '%s': %v", cur.Code, err)
	}
	
	log.Printf("Currency updated successfully: %s", cur.Code)
	return cur, nil
}

func (s *service) DeleteCurrency(code string) error {
	if code == "" {
		return fmt.Errorf("currency code cannot be empty")
	}
	
	cur, _ := s.GetCurrency(code)
	if cur == nil {
		return fmt.Errorf("cannot delete - currency '%s' not found", code)
	}
	
	if err := s.repo.DeleteCurrency(code); err != nil {
		return fmt.Errorf("failed to delete currency '%s': %v", code, err)
	}
	
	log.Printf("Currency deleted successfully: %s", code)
	return nil
}

func (s *service) ListConversions() ([]*model.Conversion, error) {
	conversions := s.repo.GetConversions()
	log.Printf("Retrieved %d conversion records", len(conversions))
	return conversions, nil
}

func (s *service) CreateConversion(amount float64, fromCode, toCode string) (*model.Conversion, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("conversion amount must be greater than zero")
	}
	if fromCode == "" || toCode == "" {
		return nil, fmt.Errorf("source and target currency codes are required")
	}

	curs := s.repo.GetCurrencies()
	from, ok1 := curs[fromCode]
	to, ok2 := curs[toCode]
	
	if !ok1 {
		return nil, fmt.Errorf("source currency '%s' not found", fromCode)
	}
	if !ok2 {
		return nil, fmt.Errorf("target currency '%s' not found", toCode)
	}
	if from.Rate <= 0 || to.Rate <= 0 {
		return nil, fmt.Errorf("invalid exchange rates - both must be positive values")
	}
	
	result := amount * (from.Rate / to.Rate)
	conv := model.NewConversion(amount, from, to, result)

	if err := s.repo.Store(conv); err != nil {
		return nil, fmt.Errorf("failed to save conversion: %v", err)
	}
	
	log.Printf("Conversion completed: %.2f %s → %.2f %s", amount, fromCode, result, toCode)
	return conv, nil
}