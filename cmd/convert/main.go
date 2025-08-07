package main

import (
	"learnpack/src/currency-converter/internal/service"
	"time"
)

func main() {
	service.InitService()

	time.Sleep(3 * time.Second)

	service.StopService()
}
