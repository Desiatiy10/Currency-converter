package main

import (
	"context"
	"fmt"
	"learnpack/src/currency-converter/internal/model"
	"learnpack/src/currency-converter/repository"
	"learnpack/src/currency-converter/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Родительский контекст и отложенная остановка всех горутин
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	//Горутина для ловли сигнала.
	//После передачи сигнала в sig вызывает cancel.
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

	//Временный код для теста.
	printData()
	usd := model.NewCurrency("USD", 1.0, "US Dollar", "$")
	eur := model.NewCurrency("EUR", 1.2, "Euro", "€")
	repository.Store(usd)
	repository.Store(eur)
	amount := 100.0
	result := amount * (eur.Rate / usd.Rate)
	conv := model.NewConversion(amount, usd, eur, result)
	repository.Store(conv)

	service.InitService(ctx)

	<-ctx.Done()
	fmt.Println("Завершаем программу...")
}

// Тестовая функция для вывода содежимого мапы и слайса
func printData() {
	currencies := repository.GetCurrencies()
	fmt.Println("\nВалюты:")
	for _, c := range currencies {
		fmt.Printf("  %s (%s) = %.2f %s\n   **********************\n", c.Code, c.Name, c.Rate, c.Symbol)
	}

	conversions := repository.GetConversions()
	fmt.Println("Конверсии: ")
	for _, conv := range conversions {
		fmt.Printf("  %.2f %s -> %.2f %s\n   **********************\n", conv.Amount, conv.From.Code, conv.Result, conv.To.Code)
	}
}
