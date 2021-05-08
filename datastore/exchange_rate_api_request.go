package datastore

import "log"

var currencies map[string]int

// LogAPIRequest logs all API request for future analytics
func LogAPIRequest(ip, from, to, date string) {
	if len(currencies) == 0 {
		currencies = GetCurrencies()
	}

	q := `INSERT INTO exchange_rate_api_requests
(ip_address, base_currency_id, converted_currency_id, date)
VALUES(?, ?, ?, ?)`

	_, err := Db.Exec(q, ip, currencies[from], currencies[to], date)

	if err != nil {
		log.Fatal(err)
	}
}
