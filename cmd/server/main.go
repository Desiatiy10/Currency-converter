package main

import (
	_ "currency-converter/docs"
	router "currency-converter/internal/app"
	grpc_server "currency-converter/internal/grpc"
	"currency-converter/proto"
	"currency-converter/repository"
	"currency-converter/service"
	"log"
	"net"
	"time"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// @title Currency Converter API
// @version 1.0
// @description REST API для управления валютами и конвертациями
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Println("Получен сигнал остановки: ", sig)
		cancel()
	}()

	if err := repository.LoadCurrenciesFromFile(); err != nil {
		fmt.Println("error загрузки валют: ", err)
	}
	if err := repository.LoadConversionsFromFile(); err != nil {
		fmt.Println("error загрузки конверсий:", err)
	}

	service.InitService(ctx)

	//Запуск REST
	srv := router.New(":8080")
	go func() {
		if err := srv.Start(); err != nil {
			fmt.Println("server error:", err)
			cancel()
		}
	}()

	//Запуск gRPC
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("error to listen %v", err)
		}
		grpcServer := grpc.NewServer()
		proto.RegisterCurrencyServiceServer(grpcServer, &grpc_server.CurrencyServer{})
		proto.RegisterConversionServiceServer(grpcServer, &grpc_server.ConversionServer{})
		fmt.Println("gRPC сервер запущен на: 9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("error run server gRPC %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Выключаем сервер...")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelShutdown()
	if err := srv.Stop(shutdownCtx); err != nil {
		fmt.Println("stopping server error:", err)
	}

	fmt.Println("Завершаем программу...")
}
