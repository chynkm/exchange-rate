package redisdb

import (
	"math"
)

func getExchangeRate(baseCurrencyCode string, exchangeRates map[string]float64) map[string]float64 {
	hash := map[string]float64{}
	for currencyCode, rate := range exchangeRates {
		roundRate := rate / exchangeRates[baseCurrencyCode]
		roundRate = math.Round(roundRate*10000) / 10000
		hash[currencyCode] = roundRate
	}

	return hash
}
