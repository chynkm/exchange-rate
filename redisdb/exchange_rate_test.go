package redisdb

import (
	"reflect"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	exchangeRates := map[string]float64{
		"EUR": 1,
		"USD": 1.2082,
		"JPY": 131.6200,
	}

	var exchangeRateTable = []struct {
		in  string
		out map[string]float64
	}{
		{
			"EUR",
			map[string]float64{
				"EUR": 1,
				"USD": 1.2082,
				"JPY": 131.62,
			},
		},
		{
			"USD",
			map[string]float64{
				"EUR": 0.8277,
				"USD": 1,
				"JPY": 108.9389,
			},
		},
		{
			"JPY",
			map[string]float64{
				"EUR": 0.0076,
				"USD": 0.0092,
				"JPY": 1,
			},
		},
	}

	for _, entry := range exchangeRateTable {
		got := getExchangeRate(entry.in, exchangeRates)
		if !reflect.DeepEqual(got, entry.out) {
			t.Error("exchange rate calculation error")
		}
	}
}
