package datastore

import (
	"database/sql"
	"log"
)

var currencies map[string]int

// LogAPIRequest logs all API request for future analytics
func LogAPIRequest(ip, from, to, date string) {
	if len(currencies) == 0 {
		currencies = GetCurrencies()
	}

	q := `INSERT INTO exchange_rate_api_requests
(ip_address, base_currency_id, converted_currency_id, date)
VALUES(?, ?, ?, ?)`

	_, err := Db.Exec(q, ip, currencies[from], newNullInt32(currencies[to]), date)

	if err != nil {
		// This isn't a critical feature. Hence, only logging
		log.Println(err)
	}
}

// newNullInt32 returns Null for empty integer
func newNullInt32(i int) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}
