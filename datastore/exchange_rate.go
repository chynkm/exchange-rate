package datastore

import (
	"log"
)

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

// GetExchangeRates returns a hash of converted currency code and their rate
// for a given date. If the "date" data is missing fetch from the last
// available date in the DB
func GetExchangeRates(date string) map[string]float64 {
	exchangeRates := map[string]float64{}
	if !exchangeRateDataExists(date) {
		date = getPreviousExchangeRateDate(date)
	}

	q := `SELECT code converted_code, rate FROM exchange_rates er
JOIN currencies c ON c.id = converted_currency_id
WHERE date = ?`

	rows, err := Db.Query(q, date)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var rate float64
		rows.Scan(&code, &rate)

		exchangeRates[code] = rate
	}

	return exchangeRates
}

// exchangeRateDataExists Check if exchange rate data exists or not
// for a particular date
func exchangeRateDataExists(date string) bool {
	var count int
	q := `SELECT COUNT(id) FROM exchange_rates WHERE date = ?`

	err := Db.QueryRow(q, date).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		return false
	}

	return true
}

// getPreviousExchangeRateDate returns the previous date containing the
// exchange rate data in the DB
func getPreviousExchangeRateDate(date string) string {
	q := `SELECT date FROM exchange_rates
WHERE date < ? ORDER BY date DESC LIMIT 1`

	err := Db.QueryRow(q, date).Scan(&date)
	if err != nil {
		log.Fatal(err)
	}

	return date
}
