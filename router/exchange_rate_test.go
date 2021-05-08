package router

import (
	"reflect"
	"testing"
	"time"

	"github.com/chynkm/ratesdb/currencystore"
	"github.com/chynkm/ratesdb/redisdb"
)

func TestGetExchangeRateQueryValidation(t *testing.T) {
	currencies := map[string]int{
		"EUR": 1,
		"USD": 2,
		"JPY": 3,
	}

	time := time.Now()
	oldDate := "1999-01-03"
	futureDate := time.AddDate(0, 0, 1).Format(currencystore.DateLayout)
	validDate1 := time.Format(currencystore.DateLayout)
	validDate2 := time.AddDate(0, 0, -10).Format(currencystore.DateLayout)

	var exchangeRateTable = []struct {
		in  map[string][]string
		out *validationError
	}{
		{
			map[string][]string{},
			&validationError{false, exchangeRateErr["from_missing"]},
		},
		{
			map[string][]string{"to": {"USD"}},
			&validationError{false, exchangeRateErr["from_missing"]},
		},
		{
			map[string][]string{"from": {"EUR"}},
			&validationError{false, exchangeRateErr["to_missing"]},
		},
		{
			map[string][]string{"from": {"EUR", "INR"}, "to": {"USD"}},
			&validationError{false, exchangeRateErr["only_one_from"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD", "INR"}},
			&validationError{false, exchangeRateErr["only_one_to"]},
		},
		{
			map[string][]string{"from": {"INR"}, "to": {"USD"}},
			&validationError{false, exchangeRateErr["unsupported_from"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}},
			&validationError{false, exchangeRateErr["unsupported_to"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {}},
			&validationError{false, exchangeRateErr["date_missing"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {"a", "b"}},
			&validationError{false, exchangeRateErr["only_one_date"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {"a"}},
			&validationError{false, exchangeRateErr["invalid_date"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {"25-04-2021"}},
			&validationError{false, exchangeRateErr["invalid_date"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {oldDate}},
			&validationError{false, exchangeRateErr["oldest_date"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"INR"}, "date": {futureDate}},
			&validationError{false, exchangeRateErr["future_date"]},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD"}},
			&validationError{true, ""},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD"}, "date": {validDate1}},
			&validationError{true, ""},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD"}, "date": {validDate2}},
			&validationError{true, ""},
		},
	}

	for _, row := range exchangeRateTable {
		got := validateGetExchangeRate(currencies, row.in)
		if !reflect.DeepEqual(got, row.out) {
			t.Error("validation error: " + row.out.message)
		}
	}
}

func TestExtractGetExchangeRateQueryParams(t *testing.T) {
	var queryParamsTable = []struct {
		in  map[string][]string
		out map[string]string
	}{
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD"}, "date": {"2021-04-25"}},
			map[string]string{"from": "EUR", "to": "USD", "date": "2021-04-25"},
		},
		{
			map[string][]string{"from": {"EUR"}, "to": {"USD"}},
			map[string]string{"from": "EUR", "to": "USD", "date": redisdb.LatestDate},
		},
	}

	for _, row := range queryParamsTable {
		date, from, to := extractGetExchangeRateQueryParams(row.in)
		if date != row.out["date"] {
			t.Error("date error")
		}
		if from != row.out["from"] {
			t.Error("from error")
		}
		if to != row.out["to"] {
			t.Error("to error")
		}
	}
}
