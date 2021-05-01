package currencystore

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

// downloadCsv the CSV file and save it to /tmp
func downloadCsv(url string) error {
	deleteCurrencyFiles()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(csvZipFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/bin/unzip", csvZipFile)
	cmd.Dir = "/tmp"
	return cmd.Run()
}

// deleteCurrencyFiles: remove existing files before downloading
func deleteCurrencyFiles() error {
	err := os.Remove(csvZipFile)
	if err != nil {
		return err
	}

	err = os.Remove(csvFile)
	return err
}
