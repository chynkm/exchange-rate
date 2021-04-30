package datastore

import (
	"testing"
)

func TestFormatDate(t *testing.T) {
	got := getDateFromString("25 April 2021")
	want := "2021-04-25"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestGetCurrencyRateFromCsv(t *testing.T) {
	input := [][]string{
		{"date", "USD", "JPY", ""},
		{"25 April 2021", "1.2070", "131.47", ""},
	}
	gotdate, gotCurrencyCode, gotRate := getCurrencyRateFromCsv(input)
	wantdate, wantCurrency, wantRate := "2021-04-25", []string{"USD", "JPY"}, []string{"1.2070", "131.47"}

	if gotdate != wantdate {
		t.Errorf("got %s, want %s", gotdate, wantdate)
	}

	if gotCurrencyCode[0] != wantCurrency[0] {
		t.Errorf("got %s, want %s", gotCurrencyCode[0], wantCurrency[0])
	}

	if gotRate[1] != wantRate[1] {
		t.Errorf("got %s, want %s", gotRate[1], wantRate[1])
	}
}
