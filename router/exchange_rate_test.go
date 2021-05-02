package router

import "testing"

func TestFromQueryValidation(t *testing.T) {
	currencies := map[string]int{
		"EUR": 1,
		"USD": 2,
		"JPY": 3,
	}

	q := map[string][]string{
		"from": {"EUR"},
		"to":   {"USD"},
	}

	got, _ := validateGetExchangeRate(currencies, q)
	if !got {
		t.Errorf("query string validation failed")
	}
}
