package main

import (
	"log"

	"github.com/ekuzm/rocket-factory/order/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
