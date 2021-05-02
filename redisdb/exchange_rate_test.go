package redisdb

import (
	"reflect"
	"testing"
)

func TestExchangeRateFormat(t *testing.T) {
	currencies := map[string]int{
		"EUR": 1,
		"USD": 2,
		"JPY": 3,
	}
	exchangeRates := map[string]float64{
		"EUR": 1,
		"USD": 1.2082,
		"JPY": 131.6200,
	}

	got := createExchangeRateHash("2021-04-25", currencies, exchangeRates)
	want := map[string]map[string]float64{
		"2021-04-25:EUR": {
			"EUR": 1,
			"USD": 1.2082,
			"JPY": 131.62,
		},
		"2021-04-25:USD": {
			"EUR": 0.8277,
			"USD": 1,
			"JPY": 108.9389,
		},
		"2021-04-25:JPY": {
			"EUR": 0.0076,
			"USD": 0.0092,
			"JPY": 1,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("exchange rate hash creation failed")
	}
}

func TestCalculateExchangeRate(t *testing.T) {
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

	for _, row := range exchangeRateTable {
		got := calculateExchangeRate(row.in, exchangeRates)
		if !reflect.DeepEqual(got, row.out) {
			t.Error("exchange rate calculation error")
		}
	}
}
