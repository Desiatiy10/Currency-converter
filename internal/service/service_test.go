package service_test

import (
	"currency-converter/internal/model"
	"currency-converter/internal/service"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	currencies  map[string]*model.Currency
	conversions []*model.Conversion
	storeErr    error
	updateErr   error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		currencies:  make(map[string]*model.Currency),
		conversions: []*model.Conversion{},
	}
}

func (m *mockRepo) LoadCurrencies() error { return nil }

func (m *mockRepo) LoadConversions() error { return nil }

func (m *mockRepo) Store(e model.Entity) error {
	if m.storeErr != nil {
		return m.storeErr
	}
	switch v := e.(type) {
	case *model.Currency:
		m.currencies[v.Code] = v
	case *model.Conversion:
		m.conversions = append(m.conversions, v)
	}
	return nil
}

func (m *mockRepo) GetCurrencies() map[string]*model.Currency {
	return m.currencies
}

func (m *mockRepo) UpdateCurrency(cur *model.Currency) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.currencies[cur.Code] = cur
	return nil
}

func (m *mockRepo) GetConversions() []*model.Conversion {
	return m.conversions
}

func TestService_CreateCurrency(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.Currency
		wantErr bool
		prepare   func(repo *mockRepo)
	}{
		{
			name: "успешное создание валюты",
			input: &model.Currency{
				Code:   "USD",
				Rate:   1.0,
				Name:   "Dollar",
				Symbol: "$",
			},
			wantErr: false,
		},
		{
			name: "некорректные данные валюты",
			input: &model.Currency{
				Code:   "",
				Rate:   -1,
				Name:   "",
				Symbol: "",
			},
			wantErr: true,
		},
		{
			name: "дублирование валюты",
			input: &model.Currency{
				Code:   "RUB",
				Rate:   1.0,
				Name:   "Рус Рубль",
				Symbol: "P",
			},
			wantErr: true,
			prepare: func(repo *mockRepo) {
				repo.Store(&model.Currency{
					Code:   "RUB",
					Rate:   1.0,
					Name:   "Рус Рубль",
					Symbol: "P",
				})
			},
		},
	}

	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(repo)
			}
			cur, err := svc.CreateCurrency(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cur)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Code, cur.Code)
			}
		})
	}
}

func TestService_GetCurrency(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	repo.currencies["EUR"] = &model.Currency{
		Code:   "EUR",
		Rate:   90,
		Name:   "Евро",
		Symbol: "€",
	}

	t.Run("валидный код", func(t *testing.T) {
		cur, err := svc.GetCurrency("EUR")
		assert.NoError(t, err)
		assert.Equal(t, "EUR", cur.Code)
	})

	t.Run("пустой код", func(t *testing.T) {
		cur, err := svc.GetCurrency("")
		assert.Error(t, err)
		assert.Nil(t, cur)
	})

	t.Run("неизвестный код", func(t *testing.T) {
		cur, err := svc.GetCurrency("ABC")
		assert.Error(t, err)
		assert.Nil(t, cur)
	})
}

func TestService_UpdateCurrency(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.Currency
		prepare func(repo *mockRepo)
		wantErr bool
	}{
		{
			name: "успешное обновление",
			input: &model.Currency{
				Code:   "GBP",
				Rate:   120,
				Name:   "Фунт стерлингов",
				Symbol: "£",
			},
			prepare: func(repo *mockRepo) {
				repo.Store(&model.Currency{Code: "GBP", Rate: 100, Name: "Фунт", Symbol: "£"})
			},
			wantErr: false,
		},
		{
			name: "некорректные данные",
			input: &model.Currency{
				Code:   "BAD",
				Rate:   0,
				Name:   "",
				Symbol: "",
			},
			wantErr: true,
		},
		{
			name: "код не 3 символа",
			input: &model.Currency{
				Code:   "LONG",
				Rate:   1,
				Name:   "Test",
				Symbol: "T",
			},
			wantErr: true,
		},
		{
			name: "ошибка обновления в репозитории",
			input: &model.Currency{
				Code:   "ERR",
				Rate:   1,
				Name:   "Test",
				Symbol: "T",
			},
			prepare: func(repo *mockRepo) {
				repo.updateErr = fmt.Errorf("update failed")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepo()
			if tt.prepare != nil {
				tt.prepare(repo)
			}
			svc := service.NewService(repo, nil)

			cur, err := svc.UpdateCurrency(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cur)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Rate, cur.Rate)
			}
		})
	}
}

func TestService_CreateConversion(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	repo.currencies["USD"] = &model.Currency{Code: "USD", Rate: 1, Name: "Доллар", Symbol: "$"}
	repo.currencies["EUR"] = &model.Currency{Code: "EUR", Rate: 2, Name: "Евро", Symbol: "€"}

	t.Run("успешная конвертация", func(t *testing.T) {
		req := &model.ConversionRequest{
			From:   "USD",
			To:     "EUR",
			Amount: 10,
		}
		conv, err := svc.CreateConversion(req)
		assert.NoError(t, err)
		assert.Equal(t, 5.0, conv.Result) 

	})

	t.Run("отрицательная сумма", func(t *testing.T) {
		req := &model.ConversionRequest{From: "USD", To: "EUR", Amount: -1}
		conv, err := svc.CreateConversion(req)
		assert.Error(t, err)
		assert.Nil(t, conv)
	})

	t.Run("валюта не найдена", func(t *testing.T) {
		req := &model.ConversionRequest{From: "XXX", To: "EUR", Amount: 10}
		conv, err := svc.CreateConversion(req)
		assert.Error(t, err)
		assert.Nil(t, conv)
	})

	t.Run("одинаковые валюты", func(t *testing.T) {
		req := &model.ConversionRequest{From: "USD", To: "USD", Amount: 10}
		conv, err := svc.CreateConversion(req)
		assert.Error(t, err)
		assert.Nil(t, conv)
	})
}

func TestService_AddEntity(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	t.Run("успешное добавление валюты", func(t *testing.T) {
		cur := &model.Currency{
			Code:   "JPY",
			Rate:   100,
			Name:   "Йена",
			Symbol: "¥",
		}
		err := svc.AddEntity(cur)
		assert.NoError(t, err)
	})

	t.Run("nil entity", func(t *testing.T) {
		err := svc.AddEntity(nil)
		assert.Error(t, err)
	})

	t.Run("переполненный канал", func(t *testing.T) {
		for i := 0; i < 56; i++ {
			_ = svc.AddEntity(&model.Currency{
				Code:   fmt.Sprintf("C%d", i),
				Rate:   1,
				Name:   "Tmp",
				Symbol: "*",
			})
		}
		err := svc.AddEntity(&model.Currency{Code: "FULL", Rate: 1, Name: "Full", Symbol: "F"})
		assert.Error(t, err)
	})
}

func TestService_ListCurrencies(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	repo.currencies["USD"] = &model.Currency{Code: "USD", Rate: 1, Name: "Доллар", Symbol: "$"}
	repo.currencies["EUR"] = &model.Currency{Code: "EUR", Rate: 2, Name: "Евро", Symbol: "€"}

	t.Run("получение всех валют", func(t *testing.T) {
		curs, err := svc.ListCurrencies()
		assert.NoError(t, err)
		assert.Len(t, curs, 2)
		assert.Contains(t, curs, "USD")
		assert.Contains(t, curs, "EUR")
	})
}

func TestService_ListConversions(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewService(repo, nil)

	from := &model.Currency{Code: "USD", Rate: 1, Name: "Доллар", Symbol: "$"}
	to := &model.Currency{Code: "EUR", Rate: 2, Name: "Евро", Symbol: "€"}
	repo.conversions = append(repo.conversions,
		model.NewConversion(10, from, to, 5),
		model.NewConversion(20, to, from, 40),
	)

	t.Run("получение всех конверсий", func(t *testing.T) {
		convs, err := svc.ListConversions()
		assert.NoError(t, err)
		assert.Len(t, convs, 2)
		assert.Equal(t, 10.0, convs[0].Amount)
		assert.Equal(t, 20.0, convs[1].Amount)
	})
}
