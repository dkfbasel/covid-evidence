// medlinetocsv is a utility to convert data from medline format to a csv table.
// The utility will parse the content of the medline file and construct a csv
// table containing all fields that are present in the file.
//
// Medline structure is as follows:
// PMID- 32074550
// OWN - NLM
// STAT- In-Process
// LR  - 20200316
// IS  - 1881-7823 (Electronic)
// IS  - 1881-7815 (Linking)
// VI  - 14
// IP  - 1
// DP  - 2020 Mar 16
// TI  - Breakthrough: Chloroquine phosphate has shown apparent efficacy in treatment of
//       COVID-19 associated pneumonia in clinical studies.
// PG  - 72-73
// LID - 10.5582/bst.2020.01047 [doi]
// AB  - The coronavirus disease 2019 (COVID-19) virus is spreading rapidly, and scientists
//       are endeavoring to discover drugs for its efficacious treatment in China.
//       Chloroquine phosphate, an old drug for treatment of malaria, is shown to have
//       apparent efficacy and acceptable safety against COVID-19 associated pneumonia in
//       multicenter clinical trials conducted in China. The drug is recommended to be
//       included in the next version of the Guidelines for the Prevention, Diagnosis, and
//       Treatment of Pneumonia Caused by COVID-19 issued by the National Health Commission
//       of the People's Republic of China for treatment of COVID-19 infection in larger
//       populations in the future.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// chars to use to separate multiple entries from the same field from each other
const multiSeparator = " | "

func main() {

	// read the input file from the command line (excluding the program name)
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalln("please pass the filepath as first argument to the utility")
	}

	medlineFileName := args[0]

	// ty to open the file
	medlineFile, err := os.Open(medlineFileName)
	if err != nil {
		log.Fatalf("could not open input file: %s, %+v", args[0], err)
	}

	// create a sc
	extension := filepath.Ext(medlineFileName)
	csvOutputPath := strings.TrimSuffix(medlineFileName, extension)
	csvOutputPath = fmt.Sprintf("%s.csv", csvOutputPath)


	// read the file line by line
	scanner := bufio.NewScanner(medlineFile)
	scanner.Split(bufio.ScanLines)

	// initialize a line counter (for error reporting)
	lineCount := 0

	// initialize a dataset counter
	datasetCount := 0

	fieldIndex := make(map[string]int)
	rows := make([]map[string]string, 0)

	currentRow := make(map[string]string)
	previousFieldName := ""

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// empty lines are used to separate datasets from each other
		if len(line) == 0 {
			rows = append(rows, currentRow)
			currentRow = make(map[string]string)
			previousFieldName = ""
			datasetCount++
			continue
		}

		// initialize string builders for field identificators and content
		// note: content can span multiple lines
		var identifier strings.Builder
		var content strings.Builder

		runeCount := -1
		for _, char := range line {
			runeCount++

			// identifier is the first 4 chars (i.e. PMID)
			if runeCount < 4 {

				// ignore spaces for identifiers
				if char == ' ' {
					continue
				}

				identifier.WriteRune(char)
			}

			if runeCount > 5 {
				content.WriteRune(char)
			}

		}

		fieldName := identifier.String();
		fieldContent := content.String();

		append := false
		if fieldName == "" {
			fieldName = previousFieldName
			append = true
		}

		err := store(fieldIndex, currentRow, append, fieldName, fieldContent)
		if err != nil {
			log.Printf("%s line: %04d, field: %s\n", err, lineCount, fieldName)
			continue
		}

		if !append {
			previousFieldName = fieldName
		}
		
	}

	// append the last row
	rows = append(rows, currentRow)

	// close the medline file
	// nolint:errcheck
	medlineFile.Close()

	// sort the field index
	fields := make([]Field, len(fieldIndex))
	var i int
	for fieldName, ordinal := range fieldIndex {
		field := Field{
			Name: fieldName,
			Ordinal: ordinal,
		}
		fields[i] = field
		i++
	}

	// sort the fields by ordinal (as appearing in the medline input file)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Ordinal < fields[j].Ordinal
	})

	// create a csv writer
	csvOutput, err := os.Create(csvOutputPath)
	if err != nil {
		log.Fatalf("could not create output csv file: %+v\n", err)
	}
	defer csvOutput.Close()

	writer := csv.NewWriter(csvOutput)
	writer.Comma = ';'
	defer writer.Flush()

	// write the name of the fields in the header column
	header := make([]string, len(fields)) 
	for i, field := range fields {
		header[i] = field.Name
	}

	err = writer.Write(header)
	if err != nil {
		log.Printf("could not write header to csv output: %+v\n", err)
	}

	// go through all rows and add them to the output
	for _, dataset := range rows {

		record := make([]string, len(fields))

		for i, field := range fields {
			content, ok := dataset[field.Name]
			if ok {
				record[i] = content
			}
		}

		err := writer.Write(record)
		if err != nil {
			log.Printf("could not write csv output: %+v\n", err)
		}
		
	}

	log.Printf("converted %d datasets from medline to csv format", len(rows))

}

// Field is used to index all fields sorted by appearance in the medline
// input field
type Field struct {
	Name string
	Ordinal int
}

// store will store the given field information into the given row
func store(fieldIndex map[string]int, row map[string]string, append bool,
	fieldName, fieldContent string) error {

		// no field name indicates, that the content should be appended
		// to existing content
		if (append) {
			// directly append to existing content
			existingContent, ok := row[fieldName]
			if !ok {
				return fmt.Errorf("could not append existing content")
			}
			row[fieldName] = fmt.Sprintf("%s%s", existingContent, fieldContent)
			return nil
		}

		// add the field to the index list if not yet present
		if _, ok := fieldIndex[fieldName]; !ok {
			fieldIndex[fieldName] = len(fieldIndex)
		}

		// append existing content with | as separator
		existingContent, ok := row[fieldName]
		if ok {
			row[fieldName] = fmt.Sprintf("%s%s%s", existingContent, 
			multiSeparator, fieldContent)
			return nil
		} 
		
		// store the content if not content is available yet
		row[fieldName] = fieldContent
		return nil
}
