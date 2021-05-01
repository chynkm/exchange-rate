package datastore

import (
	"log"
)

var currencies map[string]int

// SaveCurrencyRates saves exchange rates to the DB
func SaveExchangeRates(date string, exchangeRates map[string]float64) error {
	getCurrencies()

	values := []interface{}{}
	sqlStr := "INSERT INTO currency_rates(base_currency_id, converted_currency_id, rate, date) VALUES"

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
	stmt, _ := Db.Prepare(sqlStr)

	_, err := stmt.Exec(values...)

	return err
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
