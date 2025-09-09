package handler

import (
	"currency-converter/internal/model"
	"currency-converter/internal/usecase"
	"currency-converter/repository"

	"encoding/json"
	"net/http"
)

// CreateCurrency godoc
// @Summary Создать валюту
// @Description Добавляет новую валюту в хранилище
// @Tags currency
// @Accept json
// @Produce json
// @Param currency body model.Currency true "Валюта"
// @Success 201 {object} model.Currency
// @Failure 400 {object} map[string]string
// @Router /currency [post]
func CreateCurrency(res http.ResponseWriter, req *http.Request) {
	var cur model.Currency
	if err := json.NewDecoder(req.Body).Decode(&cur); err != nil {
		usecase.WriteError(res, http.StatusBadRequest, "невалидный JSON")
		return
	}

	if cur.Code == "" || cur.Rate <= 0 {
		usecase.WriteError(res, http.StatusBadRequest, "код и курс - обязательные поля")
		return
	}

	repository.Store(&cur)
	usecase.WriteJson(res, http.StatusCreated, cur)
}

// ListCurrencies godoc
// @Summary Получить список валют
// @Tags currency
// @Produce json
// @Success 200 {array} model.Currency
// @Router /currencies [get]
func ListCurrencies(res http.ResponseWriter, req *http.Request) {
	data := repository.GetCurrencies()

	arr := make([]*model.Currency, 0, len(data))

	for _, v := range data {
		arr = append(arr, v)
	}
	usecase.WriteJson(res, http.StatusOK, arr)
}

// GetCurrency godoc
// @Summary Получить валюту
// @Tags currency
// @Produce json
// @Param code path string true "Код валюты"
// @Success 200 {object} model.Currency
// @Failure 404 {object} map[string]string
// @Router /currency/{code} [get]
func GetCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")

	data := repository.GetCurrencies()
	if currency, ok := data[code]; ok {
		usecase.WriteJson(res, http.StatusOK, currency)
		return
	}
	usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
}

// UpdateCurrency godoc
// @Summary Обновить валюту
// @Tags currency
// @Accept json
// @Produce json
// @Param code path string true "Код валюты"
// @Param currency body model.Currency true "Валюта"
// @Success 200 {object} model.Currency
// @Failure 400 {object} map[string]string
// @Router /currency/{code} [put]
func UpdateCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")

	var upd model.Currency
	if err := json.NewDecoder(req.Body).Decode(&upd); err != nil {
		usecase.WriteError(res, http.StatusBadRequest, "невалидный JSON")
		return
	}

	upd.Code = code
	if err := repository.UpdateCurInMap(&upd); err != nil {
		usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
		return
	}

	usecase.WriteJson(res, http.StatusOK, upd)
}

// DeleteCurrency godoc
// @Summary Удалить валюту
// @Tags currency
// @Produce json
// @Param code path string true "Код валюты"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /currency/{code} [delete]
func DeleteCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")

	if err := repository.DeleteCurFromMap(code); err != nil {
		usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
		return
	}

	usecase.WriteJson(res, http.StatusOK, map[string]string{"статус:": "удалено"})
}

// CreateConversion godoc
// @Summary Создать конвертацию
// @Description Конвертирует валюту и сохраняет результат
// @Tags conversion
// @Accept json
// @Produce json
// @Param request body model.ConversionRequest true "Запрос на конвертацию"
// @Success 201 {object} model.Conversion
// @Failure 400 {object} map[string]string
// @Router /conversion [post]
func CreateConversion(res http.ResponseWriter, req *http.Request) {
	var conv struct {
		Amount float64 `json:"amount"`
		From   string  `json:"from"`
		To     string  `json:"to"`
	}
	if err := json.NewDecoder(req.Body).Decode(&conv); err != nil {
		usecase.WriteError(res, http.StatusBadRequest, "невалидный JSON")
		return
	}

	curs := repository.GetCurrencies()
	from, ok1 := curs[conv.From]
	to, ok2 := curs[conv.To]
	if !ok1 || !ok2 {
		usecase.WriteError(res, http.StatusBadRequest, "не найдена исходная или целевая валюта")
		return
	}

	if from.Rate <= 0 || to.Rate <= 0 {
		usecase.WriteError(res, http.StatusBadRequest, "курсы дллжны быть положительными")
		return
	}
	result := conv.Amount * (to.Rate / from.Rate)
	conversion := model.NewConversion(conv.Amount, from, to, result)
	repository.Store(conversion)

	usecase.WriteJson(res, http.StatusCreated, conversion)
}

// ListConversions godoc
// @Summary Получить список конвертаций
// @Tags conversion
// @Produce json
// @Success 200 {array} model.Conversion
// @Router /conversions [get]
func ListConversions(res http.ResponseWriter, req *http.Request) {
	data := repository.GetConversions()
	usecase.WriteJson(res, http.StatusOK, data)
}
