package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "currency-converter/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewCurrencyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Тесты для валют
	usd, err := client.CreateCurrency(ctx, &pb.CreateCurrencyRequest{
		Currency: &pb.Currency{
			Code:   "USD",
			Rate:   1.0,
			Name:   "US Dollar",
			Symbol: "$",
		},
	})
	if err != nil {
		log.Fatalf("error CreateCurrency: %v", err)
	}
	fmt.Println("Создана валюта:", usd)

	got, err := client.GetCurrency(ctx, &pb.Currency{Code: "USD"})
	if err != nil {
		log.Fatalf("error GetCurrency: %v", err)
	}
	fmt.Println("Получена валюта:", got)

	updated, err := client.UpdateCurrency(ctx, &pb.Currency{
		Code:   "USD",
		Rate:   13121991,
		Name:   "(-_-)",
		Symbol: "$",
	})
	if err != nil {
		log.Fatalf("error UpdateCurrency: %v", err)
	}
	fmt.Println("Обновлённая валюта:", updated)

	list, err := client.ListCurrencies(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error ListCurrencies: %v", err)
	}
	fmt.Println("Список валют:")
	for _, c := range list.Currencies {
		fmt.Printf("- %s (%s): %.2f\n", c.Code, c.Name, c.Rate)
	}

	_, err = client.DeleteCurrency(ctx, &pb.Currency{Code: "USD"})
	if err != nil {
		log.Fatalf("error DeleteCurrency: %v", err)
	}
	fmt.Println("USD удалена")

	list2, _ := client.ListCurrencies(ctx, &emptypb.Empty{})
	fmt.Println("Список валют после удаления:", list2.Currencies)

	// Тесты для конверсий
	convClient := pb.NewConversionServiceClient(conn)

	_, _ = client.CreateCurrency(ctx, &pb.CreateCurrencyRequest{
		Currency: &pb.Currency{
			Code:   "EUR",
			Rate:   0.95,
			Name:   "Euro",
			Symbol: "€",
		},
	})

	conv, err := convClient.CreateConversion(ctx, &pb.CreateConversionRequest{
		Amount: 100,
		From:   "CURB",
		To:     "CURC",
	})
	if err != nil {
		log.Fatalf("error CreateConversion: %v", err)
	}
	fmt.Printf("Конвертация: %.2f %s = %.2f %s\n",
		conv.Amount, conv.From.Code, conv.Result, conv.To.Code)

	convList, err := convClient.ListConversions(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error ListConversions: %v", err)
	}
	fmt.Println("История конверсий:")
	for _, c := range convList.Conversions {
		fmt.Printf("- %.2f %s -> %.2f %s\n", c.Amount, c.From.Code, c.Result, c.To.Code)
	}
}
