package main

import (
	"log"

	"github.com/yokoe/yagisan/internal/app/yagisan"
)

func main() {
	if err := yagisan.Run(); err != nil {
		log.Fatalf("Error: %+v\n", err)
	}
}
