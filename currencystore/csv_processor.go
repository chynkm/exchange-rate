package currencystore

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chynkm/ratesdb/datastore"
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

func GetBulkExchangeRateFromCsv() {
	currencyRates := openAndReadFile("/tmp/eurofxref-hist.csv")
	currencyCodes := currencyRates[0][1:]
	currencies := datastore.GetCurrencies()

	for i := 1; i < len(currencyRates); i++ {
		date, rates := currencyRates[i][0], currencyRates[i][1:]

		exchangeRates := map[string]float64{"EUR": 1}
		var err error

		for i, currencyCode := range currencyCodes {
			currencyCode = strings.TrimSpace(currencyCode)
			if strings.TrimSpace(rates[i]) == "" {
				exchangeRates[currencyCode] = 0
			} else {
				exchangeRates[currencyCode], err = strconv.ParseFloat(strings.TrimSpace(rates[i]), 64)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		values := []interface{}{}
		sqlStr := "INSERT INTO exchange_rates(base_currency_id, converted_currency_id, rate, date) VALUES"

		for currencyCode, rate := range exchangeRates {
			sqlStr += "(?, ?, ?, ?),"
			values = append(
				values,
				currencies["EUR"],
				currencies[currencyCode],
				rate,
				date,
			)
		}

		sqlStr = sqlStr[0 : len(sqlStr)-1]
		stmt, _ := datastore.Db.Prepare(sqlStr)

		stmt.Exec(values...)
	}
}
