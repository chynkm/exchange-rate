package datastore

import (
	"log"
	"strings"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
)

const currencyCsvUrl = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref.zip"

var currencies map[string]int

// SaveCurrencyRates saves rates to the Db
func SaveCurrencyRates() error {
	getCurrencies()
	currencystore.DownloadCsv(currencyCsvUrl)
	currencyRates := currencystore.OpenAndReadFile(currencystore.CsvFile)
	date, currencyCodes, rates := getCurrencyRateFromCsv(currencyRates)

	sqlStr := "INSERT INTO currency_rates(base_currency_id, converted_currency_id, rate, date) VALUES"
	values := []interface{}{}

	for i, currencyCode := range currencyCodes {
		sqlStr += "(?, ?, ?, ?),"
		currencyCode = strings.TrimSpace(currencyCode)
		values = append(values, currencies["EUR"], currencies[currencyCode], rates[i], date)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-1]
	stmt, _ := Db.Prepare(sqlStr)

	_, err := stmt.Exec(values...)

	return err
}

func getCurrencyRateFromCsv(currencyRates [][]string) (string, []string, []string) {
	currencyCodes := currencyRates[0][1 : len(currencyRates[0])-1]
	date, rates := currencyRates[1][0], currencyRates[1][1:len(currencyRates[0])-1]

	date = getDateFromString(date)
	return date, currencyCodes, rates
}

func getDateFromString(dt string) string {
	newdate, err := time.Parse("02 January 2006", dt)
	if err != nil {
		log.Fatal(err)
	}

	return newdate.Format("2006-01-02")
}

func getCurrencies() {
	currencies = make(map[string]int)
	q := `SELECT id, code FROM currencies`

	rows, err := Db.Query(q)

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
