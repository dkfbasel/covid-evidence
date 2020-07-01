package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// loadFromFile will load the newest export from the export files
func loadFromFile(fileName string, index map[string]int) ([]map[string]string, error) {

	file, err := os.Open(fmt.Sprintf("%s.csv", fileName))
	if err != nil {
		return nil, fmt.Errorf("could not open input file: %w", err)
	}
	defer file.Close() // nolint:errcheck
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// initialize a return slice for all records
	records := []map[string]string{}

	line := 0
	for {
		line++
		row, err := reader.Read()

		// stop when we reach the end of the record
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("could not read csv record: %w", err)
		}

		// skip the header line
		if line == 1 {
			continue
		}

		record := make(map[string]string)
		for key, colIndex := range index {
			record[key] = row[colIndex]
		}

		// define the action to perform on the record
		if record["action"] == "" {
			record["action"] = "evaluate"
		}

		records = append(records, record)

	}

	return records, nil

}

func fromCsv(fieldName string, row []string) string {
	index := -1
	for i, field := range fieldMap {
		if field.Name == fieldName {
			index = i
		}
	}

	if index == -1 {
		return ""
	}

	return row[index]
}
