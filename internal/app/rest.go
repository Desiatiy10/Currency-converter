package app

import (
	"context"
	"currency-converter/internal/handler"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	httpServer *http.Server
	curHandler *handler.CurrencyHandler
	convHandler *handler.ConversionHandler
}

func New(addr string, curHand *handler.CurrencyHandler, convHand *handler.ConversionHandler) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /currency", curHand.CreateCurrency)
	mux.HandleFunc("GET /currency/{code}", curHand.GetCurrency)
	mux.HandleFunc("GET /currencies", curHand.ListCurrencies)
	mux.HandleFunc("PUT /currency/{code}", curHand.UpdateCurrency)
	mux.HandleFunc("DELETE /currency/{code}", curHand.DeleteCurrency)

	mux.HandleFunc("POST /conversion", convHand.CreateConversion)
	mux.HandleFunc("GET /conversions", convHand.ListConversions)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		curHandler: curHand,
		convHandler: convHand,
	}
}

func (s *Server) Start() error {
	log.Println("REST server starting on: ", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Initiating server shutdown...")
	return s.httpServer.Shutdown(ctx)
}
