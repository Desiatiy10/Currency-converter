package handler

import (
	"currency-converter/internal/httputil"
	"currency-converter/internal/model"
	"currency-converter/internal/service"

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
// @Summary Get currencies list
// @Tags currency
// @Produce json
// @Success 200 {array} model.Currency
// @Failure 500 {object} map[string]string
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
// @Summary Get currency details
// @Tags currency
// @Produce json
// @Param code path string true "Currency code (e.g., USD, EUR)"
// @Success 200 {object} model.Currency
// @Failure 404 {object} map[string]string
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
// @Summary Update currency
// @Tags currency
// @Accept json
// @Produce json
// @Param code path string true "Currency code to update"
// @Param currency body model.Currency true "Updated currency data"
// @Success 200 {object} model.Currency
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
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
	}
	
	cur.Code = code
	if update, err := h.svc.UpdateCurrency(&cur); err == nil {
		httputil.WriteJson(res, http.StatusOK, update)
		return
	}
	
	httputil.WriteError(res, http.StatusNotFound, "Currency not found: "+code)
}

// DeleteCurrency godoc
// @Summary Delete currency
// @Tags currency
// @Produce json
// @Param code path string true "Currency code to delete"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /currency/{code} [delete]
func (h *CurrencyHandler) DeleteCurrency(res http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")
	if code == "" {
		httputil.WriteError(res, http.StatusBadRequest, "Currency code is required")
		return
	}

	if err := h.svc.DeleteCurrency(code); err != nil {
		httputil.WriteError(res, http.StatusNotFound, "Currency not found: "+code)
		return
	}
	
	httputil.WriteJson(res, http.StatusOK, map[string]string{
		"status": "success", 
		"message": "Currency deleted successfully: " + code,
	})
}

type ConversionHandler struct {
	svc service.Service
}

func NewConversionHandler(svc service.Service) *ConversionHandler {
	return &ConversionHandler{svc: svc}
}

// CreateConversion godoc
// @Summary Create conversion
// @Description Converts currency and saves the result
// @Tags conversion
// @Accept json
// @Produce json
// @Param request body model.ConversionRequest true "Conversion request"
// @Success 201 {object} model.Conversion
// @Failure 400 {object} map[string]string
// @Router /conversion [post]
func (h *ConversionHandler) CreateConversion(res http.ResponseWriter, req *http.Request) {
	var convReq model.ConversionRequest
	if err := httputil.ReadJson(*req, &convReq); err != nil {
		httputil.WriteError(res, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	
	conv, err := h.svc.CreateConversion(convReq.Amount, convReq.From, convReq.To)
	if err != nil {
		httputil.WriteError(res, http.StatusBadRequest, "Conversion failed: "+err.Error())
		return
	}
	
	httputil.WriteJson(res, http.StatusCreated, conv)
}

// ListConversions godoc
// @Summary Get conversions history
// @Tags conversion
// @Produce json
// @Success 200 {array} model.Conversion
// @Failure 500 {object} map[string]string
// @Router /conversions [get]
func (h *ConversionHandler) ListConversions(res http.ResponseWriter, req *http.Request) {
	data, err := h.svc.ListConversions()
	if err != nil {
		httputil.WriteError(res, http.StatusInternalServerError, "Failed to retrieve conversion history")
		return
	}
	httputil.WriteJson(res, http.StatusOK, data)
}