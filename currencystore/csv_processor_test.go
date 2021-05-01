package currencystore

import (
	"reflect"
	"testing"
)

func TestFormatDate(t *testing.T) {
	got := getDateFromString("25 April 2021")
	want := "2021-04-25"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestGetExchangeRateFromCsv(t *testing.T) {
	input := [][]string{
		{"date", " USD", " JPY", ""},
		{"25 April 2021", " 1.2070", " 131.47", ""},
	}
	gotdate, gotExchangeRates := getExchangeRateFromCsv(input)
	wantdate, wantExchangeRates := "2021-04-25", map[string]float64{
		"EUR": 1,
		"USD": 1.2070,
		"JPY": 131.47,
	}

	if gotdate != wantdate {
		t.Errorf("got %s, want %s", gotdate, wantdate)
	}

	if !reflect.DeepEqual(gotExchangeRates, wantExchangeRates) {
		t.Error("got exchangerates are different from want exchangerates")
	}
}
