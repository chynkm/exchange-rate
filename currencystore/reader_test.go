package currencystore

import (
	"bytes"
	"testing"
)

func TestReadCsvFile(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("fake, csv, file")
	_, err := readCsvFile(&buffer)

	if err != nil {
		t.Error("Failed to read csv data")
	}
}

func TestReadCsvFileContent(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString(`date, USD, JPY
28 April 2021, 1.2070, 131.47
`)
	got, _ := readCsvFile(&buffer)

	want := [][]string{
		{"date", "USD", "JPY"},
		{"28 April 2021", "1.2070", "131.47"},
	}

	if got[0][0] != want[0][0] {
		t.Errorf("got %s, want %s", got[0][0], want[0][0])
	}
}
