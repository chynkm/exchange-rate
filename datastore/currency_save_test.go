package datastore

import "testing"

func TestFormatDate(t *testing.T) {
	got := getDateFromString("25 April 2021")
	want := "2021-04-25"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
