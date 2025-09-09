package main

import (
	_ "currency-converter/docs"
	"currency-converter/internal/app"
	"currency-converter/repository"
	"currency-converter/service"
	"time"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// @title Currency Converter API
// @version 1.0
// @description REST API для управления валютами и конвертациями
// @host localhost:8080
// @BasePath /
func main() {
	//Родительский контекст и отложенная остановка всех горутин
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Println("Получен сигнал остановки: ", sig)
		cancel()
	}()

	if err := repository.LoadCurrenciesFromFile(); err != nil {
		fmt.Println("ошибка загрузки валют: ", err)
	}
	if err := repository.LoadConversionsFromFile(); err != nil {
		fmt.Println("ошибка загрузки конвертаций:", err)
	}

	service.InitService(ctx)

	srv := router.New(":8080")
	go func() {
		if err := srv.Start(); err != nil {
			fmt.Println("ошибка сервера:", err)
			cancel()
		}
	}()

	<-ctx.Done()

	// Graceful shutdown
	fmt.Println("Выключаем сервер...")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelShutdown()
	if err := srv.Stop(shutdownCtx); err != nil {
		fmt.Println("ошибка при завершении сервера:", err)
	}

	fmt.Println("Завершаем программу...")
}
