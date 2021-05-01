package datastore

import "log"

// SaveCurrencyRates saves exchange rates to the DB
func SaveExchangeRates(date string, exchangeRates map[string]float64) error {
	currencies := GetCurrencies()

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
	stmt, _ := Db.Prepare(sqlStr)

	_, err := stmt.Exec(values...)

	return err
}

// GetCurrencies generates a hash of currency code and its DB id value
func GetCurrencies() map[string]int {
	currencies := map[string]int{}

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

	return currencies
}

// GetExchangeRates returns a hash of converted currency code
// and their rate for a given date
func GetExchangeRates(date string) map[string]float64 {
	exchangeRates := map[string]float64{}

	q := `SELECT code converted_code, rate FROM exchange_rates er
JOIN currencies c ON c.id = converted_currency_id
WHERE date = ?`

	rows, err := Db.Query(q, date)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var code string
		var rate float64
		rows.Scan(&code, &rate)

		exchangeRates[code] = rate
	}

	return exchangeRates
}
