package main

import (
	"log"

	"github.com/chynkm/exchange-rate/datastore"
)

func main() {
	err := datastore.SaveCurrencyRates()
	if err != nil {
		log.Fatal(err)
	}
}
