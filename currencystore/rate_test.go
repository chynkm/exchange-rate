package currencystore

import (
	"testing"
)

func TestDownloadCsv(t *testing.T) {
	url := ""
	err := DownloadCsv(url)

	if err == nil {
		t.Errorf("Empty URL should create error")
	}
}
