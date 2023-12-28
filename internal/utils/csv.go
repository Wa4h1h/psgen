package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadAllCsv(path string, delim rune) ([][]string, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening csv error: %w", err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Comma = delim
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading all csv content error: %w", err)
	}

	return records, nil
}
