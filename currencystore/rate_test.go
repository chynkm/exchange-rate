package currencystore

import (
	"testing"
)

func TestDownloadCsv(t *testing.T) {
	url := ""
	err := downloadCsv(url)

	if err == nil {
		t.Errorf("Empty URL should create error")
	}
}
