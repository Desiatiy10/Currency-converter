package main

import (
	serv "learnpack/src/currency-converter/internal/service"
	"time"
)

func main() {
	serv.InitService()

	time.Sleep(time.Second * 1)

	serv.StopService()
}
