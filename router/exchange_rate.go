package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chynkm/ratesdb/datastore"
)

var currencies map[string]int

type ExchangeRate struct {
	rate float64
}

type ValidationError struct {
	status  int
	message string
}

func getExchangeRate(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(map[string]int{"status": 200})
}

func validateGetExchangeRate(currencies map[string]int, q map[string][]string) (bool, string) {
	from := q["from"][0]
	if _, ok := currencies[from]; !ok {
		return false, ""
	}

	return true, ""
}

// Routes holds all the routes supported by the application
func Routes() {
	currencies = datastore.GetCurrencies()
	http.HandleFunc("/v1/rates", getExchangeRate)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
