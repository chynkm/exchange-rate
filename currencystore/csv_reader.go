package currencystore

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func readCsvFile(reader io.Reader) ([][]string, error) {
	r := csv.NewReader(reader)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return lines, err
}

func openAndReadFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open file: %s", fileName)
	}

	lines, err := readCsvFile(file)
	if err != nil {
		fmt.Printf("Failed to read file: %s", fileName)
	}
	return lines
}
