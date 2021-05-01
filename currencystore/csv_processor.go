package currencystore

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// FetchExchangeRates download the latest file from central bank of Europe,
// processes the file and returns the exchange rates
func FetchExchangeRates() (string, map[string]float64) {
	downloadCsv(currencyCsvUrl)
	currencyRates := openAndReadFile(csvFile)
	return getExchangeRateFromCsv(currencyRates)
}

func getExchangeRateFromCsv(currencyRates [][]string) (string, map[string]float64) {
	currencyCodes := currencyRates[0][1 : len(currencyRates[0])-1]
	date, rates := currencyRates[1][0], currencyRates[1][1:len(currencyRates[0])-1]

	date = getDateFromString(date)

	exchangeRates := map[string]float64{"EUR": 1}
	var err error

	for i, currencyCode := range currencyCodes {
		currencyCode = strings.TrimSpace(currencyCode)
		exchangeRates[currencyCode], err = strconv.ParseFloat(strings.TrimSpace(rates[i]), 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	return date, exchangeRates
}

func getDateFromString(dt string) string {
	newdate, err := time.Parse("02 January 2006", dt)
	if err != nil {
		log.Fatal(err)
	}

	return newdate.Format("2006-01-02")
}
