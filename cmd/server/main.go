package main

import (
	_ "currency-converter/docs"
	"currency-converter/internal/app"
	"currency-converter/internal/handler"
	"currency-converter/internal/repository"
	"currency-converter/internal/service"
	"currency-converter/proto"

	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// @title Currency Converter API
// @version 1.0
// @description REST API для управления валютами и конвертациями
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Println("Received shutdown signal: ", sig)
		cancel()
	}()

	// Repository
	repo := repository.NewRepository()
	if err := repo.LoadCurrencies(); err != nil {
		fmt.Println("Failed to load currency data: ", err)
	}
	if err := repo.LoadConversions(); err != nil {
		fmt.Println("Failed to load conversion data:", err)
	}

	//Service
	srvc := service.InitService(ctx, repo)

	//Handlers
	curHandler := handler.NewCurrencyHandler(srvc)
	convHandler := handler.NewConversionHandler(srvc)

	//Запуск REST
	server := app.New(":8080", curHandler, convHandler)
	go func() {
		if err := server.Start(); err != nil {
			fmt.Println("REST API server error: ", err)
			cancel()
		}
	}()

	//Запуск gRPC
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()

		proto.RegisterCurrencyServiceServer(grpcServer,
			app.NewCurrencyServer(srvc))
		proto.RegisterConversionServiceServer(grpcServer,
			app.NewConversionServer(srvc))

		log.Println("gRPC server starting on port: 9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(), 2*time.Second)
	defer cancelShutdown()
	if err := server.Stop(shutdownCtx); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}

	fmt.Println("Application terminated successfully")
}
