package redisdb

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/datastore"
)

// SaveExchangeRates to Redis
// Insert data of previous days exchange rate when a day is missing
// Previous day is always present since we are fetching it from the DB
func SaveExchangeRates() {
	rdb := Rdbpool.Get()
	defer rdb.Close()

	dbCurrencies := datastore.GetCurrencies()

	sDate := datastore.GetOldestExchangeRateDate(days)
	startdate, err := time.Parse("2006-01-02", sDate)
	endDate := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	dates := generateDates(startdate, endDate)
	fmt.Println(dates)
	os.Exit(1)

	for i := 0; i < len(dates); i++ {
		exchangeRates := datastore.GetExchangeRates(dates[i])

		dailyExchangeRates := createExchangeRateHash(
			dates[i],
			dbCurrencies,
			exchangeRates,
		)

		for key, dailyExchangeRate := range dailyExchangeRates {
			redisExchangeRates := []interface{}{}
			redisExchangeRates = append(redisExchangeRates, key)
			for code, rate := range dailyExchangeRate {
				redisExchangeRates = append(redisExchangeRates, code, rate)
			}

			_, err := rdb.Do("HMSET", redisExchangeRates...)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// createExchangeRateHash generates hash for exchange rates of different
// currencies for the given date
func createExchangeRateHash(
	date string,
	currencies map[string]int,
	exchangeRates map[string]float64,
) map[string]map[string]float64 {
	exchangeRateHash := map[string]map[string]float64{}

	for currencyCode, _ := range currencies {
		exchangeRateHash[date+":"+currencyCode] = getExchangeRate(
			currencyCode,
			exchangeRates,
		)
	}

	return exchangeRateHash
}

// getExchangeRate for a base currency
func getExchangeRate(
	baseCurrencyCode string,
	exchangeRates map[string]float64,
) map[string]float64 {
	if baseCurrencyCode == euro {
		return exchangeRates
	}

	baseCurrencyExchangeRate := map[string]float64{}
	for currencyCode, rate := range exchangeRates {
		roundRate := rate / exchangeRates[baseCurrencyCode]
		roundRate = math.Round(roundRate*10000) / 10000
		baseCurrencyExchangeRate[currencyCode] = roundRate
	}

	return baseCurrencyExchangeRate
}

// generateDates from the start date to the end date.
// It includes the start date so that corresponding DB value is present
func generateDates(startDate time.Time, endDate time.Time) []string {
	dates := []string{}
	i := 0
	for startDate.AddDate(0, 0, i).Format(currencystore.DateLayout) <= endDate.Format(currencystore.DateLayout) {
		dates = append(dates, startDate.AddDate(0, 0, i).Format(currencystore.DateLayout))
		i++
	}

	return dates
}
