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

func WriteToCsv(header []string, body [][]string, path string, delim rune) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("creating csv file error: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = delim

	if err := w.Write(header); err != nil {
		return fmt.Errorf("writing to csv error: %w", err)
	}

	for _, row := range body {
		if err := w.Write(row); err != nil {
			return fmt.Errorf("writing to csv error: %w", err)
		}
	}

	return nil
}
