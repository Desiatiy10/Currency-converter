package main

import (
	"context"
	"learnpack/src/currency-converter/service"
	"time"
)

func main() {
	//Родительский контекст и отложенная остановка всех горутин
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service.InitService(ctx)

	// usd := model.NewCurrency("USD", 1.00, "Dollar", "$")
	// rub := model.NewCurrency("RUB", 0.012484, "Рубль", "₽")

	// conv := model.NewConversion(100, usd, rub, 8010)

	// fmt.Printf(conv)

	time.Sleep(time.Second * 7)
}
