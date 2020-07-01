package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

// Parse will convert the given information to the specified data model
func Parse(inputFile string) error {

	jsonFileName := fmt.Sprintf("%s.json", inputFile)

	content, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		return fmt.Errorf("could not open json file: %w", err)
	}

	// parse the content of the json file
	var studies []json.RawMessage
	err = json.Unmarshal(content, &studies)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	fmt.Println("number of studies", len(studies))

	// initialize a csv file for the output
	csvFileName := strings.Replace(jsonFileName, ".json", ".csv", 1)
	file, err := os.Create(csvFileName)
	if err != nil {
		return fmt.Errorf("could not create csv file: %w", err)
	}
	defer file.Close() // nolint:errcheck

	// initialize a new csv writer
	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// write the header line
	rowCount := len(fieldMap)
	header := make([]string, rowCount)
	for i, field := range fieldMap {
		header[i] = field.Name
	}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("could not write header to csv output: %w", err)
	}

	// iterate through all studies in the dataset
	for _, study := range studies {

		// initialize a searchable json structure
		study := gjson.GetBytes(study, "Study")

		// initialize a new row
		row := make([]string, rowCount)

		// try to extract the field content according to the field map
		for i, field := range fieldMap {

			// skip all fields without a search path specified
			if field.Search == "" {
				continue
			}

			value := study.Get(field.Search)
			if !value.Exists() {
				// fmt.Printf("could not find field: %s\n", field.Search)
				continue
			}

			// concatenate arrays with semicolon
			// note: might be better to keep items separated in the future
			if value.IsArray() {
				values := value.Array()
				asStr := make([]string, len(values))
				for j, item := range values {
					asStr[j] = item.String()
				}
				row[i] = strings.Join(asStr, "; ")
				continue
			}

			row[i] = value.String()

		}

		// write the row to the csv file
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("could not write study to csv output: %w", err)
		}
		writer.Flush()
	}

	file.Close() // nolint:errcheck

	return nil

}
