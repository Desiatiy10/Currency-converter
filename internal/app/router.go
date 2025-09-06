package router

import (
	"currency-converter/internal/handler"
	"currency-converter/internal/usecase"

	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/currency/create", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			handler.CreateCurrency(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/currencies/get", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			handler.ListCurrencies(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/currency/get/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			handler.GetCurrency(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/currency/put/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut {
			handler.UpdateCurrency(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/currency/delete/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodDelete {
			handler.DeleteCurrency(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/conversion/create", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			handler.CreateConversion(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})

	http.HandleFunc("/conversions/get", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			handler.ListConversions(res, req)
			return
		}
		usecase.WriteError(res, http.StatusMethodNotAllowed, "метод не поддерживается")
	})
}
