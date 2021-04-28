package datastore

import (
	"fmt"
	"log"

	"github.com/chynkm/exchange-rate/currencystore"
)

const currencyCsvUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref.zip"

var currencies map[string]int

// SaveCurrencyRates saves rates to the DB
func SaveCurrencyRates() {
	currencystore.DownloadCsv(currencyCsvUrl)
	currencyRates := currencystore.OpenAndReadFile(currencystore.CsvFile)

	currencyCodes := currencyRates[0][1 : len(currencyRates[0])-1]
	date, rates := currencyRates[1][0], currencyRates[1][1:len(currencyRates[0])-1]

	fmt.Println(date)

	for i, currency := range currencyCodes {
		fmt.Println(currency + " => " + rates[i])
	}
}

func getCurrencies() {
	currencies = make(map[string]int)
	q := `SELECT id, code FROM currencies`

	rows, err := db.Query(q)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var code string
		rows.Scan(&id, &code)

		currencies[code] = id
	}
}
