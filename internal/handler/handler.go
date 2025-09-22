package handler

import (
	"currency-converter/internal/httputil"
	"currency-converter/internal/model"
	"currency-converter/internal/service"
	"strings"

	"net/http"
)

type CurrencyHandler struct {
	svc service.Service
}

func NewCurrencyHandler(svc service.Service) *CurrencyHandler {
	return &CurrencyHandler{svc: svc}
}

// CreateCurrency godoc
// @Summary Create currency
// @Description Adds a new currency to the storage
// @Tags currency
// @Accept json
// @Produce json
// @Param currency body model.Currency true "Currency data"
// @Success 201 {object} model.Currency
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /currency [post]
func (h *CurrencyHandler) CreateCurrency(res http.ResponseWriter, req *http.Request) {
	var cur model.Currency
	if err := httputil.ReadJson(*req, &cur); err != nil {
		httputil.WriteError(res, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if respCur, err := h.svc.CreateCurrency(&cur); err == nil {
		httputil.WriteJson(res, http.StatusCreated, &respCur)
		return
	}

	httputil.WriteError(res, http.StatusBadRequest, "Invalid currency data provided")
}

// ListCurrencies godoc
// @Summary Get list of all available currencies
// @Description Retrieves all currencies with current exchange rates from Central Bank of Russia
// @Tags currency
// @Produce json
// @Success 200 {object} map[string]model.Currency "Successfully retrieved currencies map"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /currencies [get]
func (h *CurrencyHandler) ListCurrencies(res http.ResponseWriter, req *http.Request) {
	data, err := h.svc.ListCurrencies()
	if err != nil {
		httputil.WriteError(res, http.StatusInternalServerError, "Failed to retrieve currencies list")
		return
	}
	httputil.WriteJson(res, http.StatusOK, data)
}

// GetCurrency godoc
// @Summary Get currency details by code
// @Description Retrieves detailed information about specific currency including exchange rate from Central Bank of Russia
// @Tags currency
// @Produce json
// @Param code path string true "Currency code (ISO 4217 format)" Example(USD)
// @Success 200 {object} model.Currency "Successfully retrieved currency details"
// @Failure 400 {object} map[string]string "Invalid currency code format"
// @Failure 404 {object} map[string]string "Currency not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /currency/{code} [get]
func (h *CurrencyHandler) GetCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")
	if code == "" {
		httputil.WriteError(res, http.StatusBadRequest, "Currency code is required")
		return
	}

	cur, err := h.svc.GetCurrency(code)
	if err != nil {
		httputil.WriteError(res, http.StatusNotFound, "Currency not found: "+code)
		return
	}
	httputil.WriteJson(res, http.StatusOK, cur)
}

// UpdateCurrency godoc
// @Summary Update currency exchange rate
// @Description Updates the exchange rate for a specific currency. Note: Normally rates are updated automatically from Central Bank of Russia
// @Tags currency
// @Accept json
// @Produce json
// @Param code path string true "Currency code to update (ISO 4217 format)" Example(USD)
// @Param currency body model.Currency true "Currency data with updated exchange rate"
// @Success 200 {object} model.Currency "Successfully updated currency"
// @Failure 400 {object} map[string]string "Invalid input data or currency code mismatch"
// @Failure 404 {object} map[string]string "Currency not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /currency/{code} [put]
func (h *CurrencyHandler) UpdateCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")
	if code == "" {
		httputil.WriteError(res, http.StatusBadRequest, "Currency code is required")
		return
	}

	var cur model.Currency
	if err := httputil.ReadJson(*req, &cur); err != nil {
		httputil.WriteError(res, http.StatusBadRequest, "Invalid JSON format")
		return
	} else if cur.Rate <= 0 {
		httputil.WriteError(res, http.StatusBadRequest, "Exchange rate must be greater than zero")
		return
	}

	cur.Code = code
	if update, err := h.svc.UpdateCurrency(&cur); err == nil {
		httputil.WriteJson(res, http.StatusOK, update)
		return
	}

	httputil.WriteError(res, http.StatusNotFound, "Currency not found: "+code)
}

type ConversionHandler struct {
	svc service.Service
}

func NewConversionHandler(svc service.Service) *ConversionHandler {
	return &ConversionHandler{svc: svc}
}

// CreateConversion godoc
// @Summary Convert currency amount
// @Description Converts amount from one currency to another using current Central Bank of Russia exchange rates and saves the conversion result
// @Tags conversion
// @Accept json
// @Produce json
// @Param request body model.ConversionRequest true "Conversion request parameters" Example({"amount": 100, "from": "USD", "to": "EUR"})
// @Success 201 {object} model.Conversion "Successfully converted currency"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Currency not found"
// @Failure 422 {object} map[string]string "Invalid conversion parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /conversion [post]
func (h *ConversionHandler) CreateConversion(res http.ResponseWriter, req *http.Request) {
	var convReq model.ConversionRequest
	if err := httputil.ReadJson(*req, &convReq); err != nil {
		httputil.WriteError(res, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	conv, err := h.svc.CreateConversion(convReq.Amount, convReq.From, convReq.To)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			httputil.WriteError(res, http.StatusNotFound, err.Error())
		case strings.Contains(err.Error(), "must be greater than zero"),
			strings.Contains(err.Error(), "cannot convert"):
			httputil.WriteError(res, http.StatusUnprocessableEntity, err.Error())
		default:
			httputil.WriteError(res, http.StatusBadRequest, err.Error())
		}
		return
	}

	httputil.WriteJson(res, http.StatusCreated, conv)
}

// ListConversions godoc
// @Summary Get conversion history
// @Description Retrieves history of all currency conversions performed
// @Tags conversion
// @Produce json
// @Success 200 {array} model.Conversion "Successfully retrieved conversion history"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /conversions [get]
func (h *ConversionHandler) ListConversions(res http.ResponseWriter, req *http.Request) {
	data, err := h.svc.ListConversions()
	if err != nil {
		httputil.WriteError(res, http.StatusInternalServerError, "Failed to retrieve conversion history")
		return
	}
	httputil.WriteJson(res, http.StatusOK, data)
}
