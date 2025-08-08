package main

import (
	"learnpack/src/currency-converter/internal/service"
	"time"
)

func main() {
	service.InitService()

	time.Sleep(time.Second)

	service.StopService()
}
