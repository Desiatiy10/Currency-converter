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
}

func New(addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /currency", handler.CreateCurrency)
	mux.HandleFunc("GET /currency/{code}", handler.GetCurrency)
	mux.HandleFunc("GET /currencies", handler.ListCurrencies)
	mux.HandleFunc("PUT /currency/{code}", handler.UpdateCurrency)
	mux.HandleFunc("DELETE /currency/{code}", handler.DeleteCurrency)

	mux.HandleFunc("POST /conversion", handler.CreateConversion)
	mux.HandleFunc("GET /conversions", handler.ListConversions)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	log.Println("Сервер запущен на: ", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Остановка сервера...")
	return s.httpServer.Shutdown(ctx)
}
