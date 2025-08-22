package main

import (
	"context"
	"fmt"
	"learnpack/src/currency-converter/service"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	service.InitService(ctx)

	<-ctx.Done()
	fmt.Println("Ожидание завершения всех процессов.")

	time.Sleep(time.Second * 1)

	fmt.Println("Приложение завершено корректно.")
}
