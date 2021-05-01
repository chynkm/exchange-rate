package redisdb

import (
	"math"
)

const euro = "EUR"

// createExchangeRateHash generates hash for exchange rates of different
// currencies for the given date
func createExchangeRateHash(
	date string,
	currencies map[string]int,
	exchangeRates map[string]float64,
) map[string]map[string]float64 {
	exchangeRateHash := map[string]map[string]float64{}

	for dbCurrencyCode, _ := range currencies {
		exchangeRateHash[date+":"+dbCurrencyCode] = getExchangeRate(
			dbCurrencyCode,
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
