package datastore

import "log"

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
