package service

import (
	"context"
	"currency-converter/internal/api/cbr"
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

	ListConversions() ([]*model.Conversion, error)
	CreateConversion(amount float64, fromCode, toCode string) (*model.Conversion, error)
}

type service struct {
	repo       repository.Repository
	entityChan chan model.Entity
	cbrClient  *cbr.CBRClient
}

func NewService(repo repository.Repository) *service {
	return &service{
		repo:       repo,
		entityChan: make(chan model.Entity, 56),
		cbrClient:  cbr.NewCBRClient(),
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

func (s *service) syncCBRData(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	semaphore := make(chan struct{}, 3)

	if err := s.loadCBRData(ctx); err != nil {
		log.Printf("Failed to load initial ЦБ РФ data: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			select {
			case semaphore <- struct{}{}:
				go func(ctx context.Context) {
					defer func() {
						<-semaphore
						if r := recover(); r != nil {
							log.Printf("Panic recovered in data ЦБ РФ %v", r)
						}
					}()
					if err := s.loadCBRData(ctx); err != nil {
						log.Printf("Failed to sync ЦБ РФ data: %v", err)
					}
				}(ctx)
			default:
				log.Println("ЦБ РФ sync skipped: too many concurrent requests")
			}
		case <-ctx.Done():
			log.Println("ЦБ РФ data sync stopped: context cancelled")
			return
		}
	}
}

func (s *service) loadCBRData(ctx context.Context) error {
	rates, err := s.cbrClient.GetDailyRates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ЦБ РФ rates: %w", err)
	}

	// ---The Russian ruble is the base currency---
	baseRates := make(map[string]*model.Currency)
	baseRates["RUB"] = &model.Currency{
		Code:   "RUB",
		Rate:   1.0,
		Name:   "Российский рубль",
		Symbol: "₽",
	}

	for code, rate := range rates.Valute {
		rates := rate.Value / rate.Nominal
		baseRates[code] = &model.Currency{
			Code:   code,
			Rate:   rates,
			Name:   rate.Name,
			Symbol: getCurrencySymbol(code),
		}
	}

	for _, currency := range baseRates {
		if err := s.AddEntity(currency); err != nil {
			log.Printf("Failed to store currency %s: %v", currency.Code, err)
		}
	}

	time.Sleep(time.Millisecond)
	log.Printf("Loaded %d currencies from ЦБ РФ", len(baseRates))
	return nil
}

func getCurrencySymbol(code string) string {
	symbols := map[string]string{
		"AUD": "A$",     // Австралийский доллар
		"AZN": "₼",      // Азербайджанский манат
		"DZD": "د.ج",    // Алжирский динар
		"GBP": "£",      // Фунт стерлингов
		"AMD": "֏",      // Армянский драм
		"BHD": ".ب.د",   // Бахрейнский динар
		"BYN": "Br",     // Белорусский рубль
		"BGN": "лв",     // Болгарский лев
		"BOB": "Bs.",    // Боливиано
		"BRL": "R$",     // Бразильский реал
		"HUF": "Ft",     // Венгерский форинт
		"VND": "₫",      // Вьетнамский донг
		"HKD": "HK$",    // Гонконгский доллар
		"GEL": "₾",      // Грузинский лари
		"DKK": "kr",     // Датская крона
		"AED": "د.إ",    // Дирхам ОАЭ
		"USD": "$",      // Доллар США
		"EUR": "€",      // Евро
		"EGP": "ج.م",    // Египетский фунт
		"INR": "₹",      // Индийская рупия
		"IDR": "Rp",     // Индонезийская рупия
		"IRR": "﷼",      // Иранский риал
		"KZT": "₸",      // Казахстанский тенге
		"CAD": "C$",     // Канадский доллар
		"QAR": "ر.ق",    // Катарский риал
		"KGS": "сом",    // Киргизский сом
		"CNY": "¥",      // Китайский юань
		"CUP": "C$",     // Кубинское песо
		"MDL": "lei",    // Молдавский лей
		"MNT": "₮",      // Монгольский тугрик
		"NGN": "₦",      // Нигерийская найра
		"NZD": "NZ$",    // Новозеландский доллар
		"NOK": "kr",     // Норвежская крона
		"OMR": "ر.ع",    // Оманский риал
		"PLN": "zł",     // Польский злотый
		"SAR": "ر.س",    // Саудовский риял
		"RON": "lei",    // Румынский лей
		"XDR": "XDR",    // СДР
		"SGD": "S$",     // Сингапурский доллар
		"TJS": "сомони", // Таджикский сомони
		"THB": "฿",      // Тайский бат
		"BDT": "৳",      // Бангладешская така
		"TRY": "₺",      // Турецкая лира
		"TMT": "m",      // Туркменский манат
		"UZS": "so'm",   // Узбекский сум
		"UAH": "₴",      // Украинская гривна
		"CZK": "Kč",     // Чешская крона
		"SEK": "kr",     // Шведская крона
		"CHF": "CHF",    // Швейцарский франк
		"ETB": "Br",     // Эфиопский быр
		"RSD": "din",    // Сербский динар
		"ZAR": "R",      // Южноафриканский рэнд
		"KRW": "₩",      // Южнокорейский вон
		"JPY": "¥",      // Японская иена
		"MMK": "K",      // Мьянманский кьят
		"RUB": "₽",      // Российский рубль
	}
	if symbol, exist := symbols[code]; exist {
		return symbol
	}
	return code
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

func InitService(ctx context.Context, repo repository.Repository, cbrClient *cbr.CBRClient) *service {
	s := NewService(repo)

	go s.processEntities(ctx)
	go s.syncCBRData(ctx)
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

	if err := s.AddEntity(cur); err != nil {
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

func (s *service) ListConversions() ([]*model.Conversion, error) {
	conversions := s.repo.GetConversions()
	log.Printf("Retrieved %d conversion records", len(conversions))
	return conversions, nil
}

func (s *service) CreateConversion(nominal float64, fromCode, toCode string) (*model.Conversion, error) {
	if nominal <= 0 {
		return nil, fmt.Errorf("conversion amount must be greater than zero")
	}
	if fromCode == "" || toCode == "" {
		return nil, fmt.Errorf("source and target currency codes are required")
	}

	curs := s.repo.GetCurrencies()
	from, ok1 := curs[fromCode]
	if !ok1 {
		return nil, fmt.Errorf("source currency '%s' not found", fromCode)
	} else if from.Rate <= 0 {
		return nil, fmt.Errorf("invalid exchange rates - both must be positive values")
	}
	to, ok2 := curs[toCode]
	if !ok2 {
		return nil, fmt.Errorf("target currency '%s' not found", toCode)
	} else if to.Rate <= 0 {
		return nil, fmt.Errorf("invalid exchange rates - both must be positive values")
	}

	nominalInRubles := nominal * from.Rate
	result := nominalInRubles / to.Rate

	conv := model.NewConversion(nominal, from, to, result)

	if err := s.AddEntity(conv); err != nil {
		return nil, fmt.Errorf("failed to save conversion: %v", err)
	}

	log.Printf("Conversion completed: %.2f %s → %.2f %s", nominal, fromCode, result, toCode)
	return conv, nil
}
