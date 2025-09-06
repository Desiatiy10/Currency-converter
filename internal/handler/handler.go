package handler

import (
	"currency-converter/internal/model"
	"currency-converter/internal/usecase"
	"currency-converter/repository"

	"encoding/json"
	"net/http"
	"strings"
)

// Получает валюту в теле запроса и отдает
// в repository.Store через URL /currency/create c методом POST.
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

// Получает копию мапы валют и отправляет ее ответом на запрос /currencies/get с методом GET.
func ListCurrencies(res http.ResponseWriter, req *http.Request) {
	data := repository.GetCurrencies()

	arr := make([]*model.Currency, 0, len(data))

	for _, v := range data {
		arr = append(arr, v)
	}
	usecase.WriteJson(res, http.StatusOK, arr)
}

// Получает код валюты из URL /currency/get/{code} c методом GET
// находит в копии мапы и отправляет ее ответом.
func GetCurrency(res http.ResponseWriter, req *http.Request) {
	codePart := strings.TrimPrefix(req.URL.Path, "/currency/get/")
	if codePart == "" {
		usecase.WriteError(res, http.StatusBadRequest, "некорректный путь")
		return
	}

	data := repository.GetCurrencies()
	cur, ok := data[codePart]
	if !ok {
		usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
		return
	}
	usecase.WriteJson(res, http.StatusOK, cur)
}

// Получает код валюты из URL /currency/get/{code} c методом PUT
// По телу запроса и модели molel.Currency находит в мапе и обновялет данные.
func UpdateCurrency(res http.ResponseWriter, req *http.Request) {
	codePart := strings.TrimPrefix(req.URL.Path, "/currency/put/")
	if codePart == "" {
		usecase.WriteError(res, http.StatusBadRequest, "некорректный путь")
		return
	}

	var upd model.Currency
	if err := json.NewDecoder(req.Body).Decode(&upd); err != nil {
		usecase.WriteJson(res, http.StatusBadRequest, "невалидный JSON")
		return
	}

	upd.Code = codePart
	if err := repository.UpdataCurInMap(&upd); err != nil {
		usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
		return
	}

	repository.Store(&upd)
	usecase.WriteJson(res, http.StatusOK, upd)
}

// Получает код валюты из URL /currency/delete/{code} c методом DELETE
// находит в мапе и удаляет данные.
func DeleteCurrency(res http.ResponseWriter, req *http.Request) {
	codePart := strings.TrimPrefix(req.URL.Path, "/currency/delete/")
	if codePart == "" {
		usecase.WriteError(res, http.StatusBadRequest, "некорректный путь")
		return
	}

	if err := repository.DeleteCurFromMap(codePart); err != nil {
		usecase.WriteError(res, http.StatusNotFound, "валюта не найдена")
		return
	}

	usecase.WriteJson(res, http.StatusOK, map[string]string{"статус:": "удалено"})
}

// Получает валюты и сумму в теле запроса по URL /conversion/create c методом PUT
// выполняет конвертацию по тестовой формуле и отправляет в repository.Store.
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

// Получает копию слайса конвертаций по запросу /conversions/get с методом GET
func ListConversions(res http.ResponseWriter, req *http.Request) {
	data := repository.GetConversions()
	usecase.WriteJson(res, http.StatusOK, data)
}
