package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"dkfbasel.ch/covid-evidence/ninox"
)

// Compare will compare the current clinicaltrials information with the
// information stored in the ninox database and list all updates accordingly
func Compare(inputFile string) error {

	// compile a list of all fields that should be checked in ninox
	fields := []string{"action", "covebasic"}
	fieldInputRowIndex := make(map[string]int)

	for i, field := range fieldMap {
		if field.Ninox != "" {
			fields = append(fields, field.Ninox)
			fieldInputRowIndex[field.Ninox] = i
		}
	}

	// first load data from the local file
	fromSource, err := loadFromFile(inputFile, fieldInputRowIndex)
	if err != nil {
		return fmt.Errorf("could not load records from file")
	}

	fmt.Printf("clinicaltrials: %d\n", len(fromSource))

	// now fetch all data from ninox
	ninoxScreeningClinicaltrials, err := ninox.FetchRecords(ninox.ClinicaltrialsURL, "")
	if err != nil {
		return fmt.Errorf("could not fetch clnicaltrials records from ninox")
	}

	ninoxBasicIncluded, err := ninox.FetchRecords(ninox.CoveBasicURL, "")
	if err != nil {
		return fmt.Errorf("could not fetch covebasic records from ninox")
	}

	ninoxBasicExcluded, err := ninox.FetchRecords(ninox.CoveBasicExlusionURL, "")
	if err != nil {
		return fmt.Errorf("could not fetch covebasic-exlude records from ninox")
	}

	fmt.Printf("from ninox, clinialtrials: %d\n", len(ninoxScreeningClinicaltrials))
	fmt.Printf("from ninox, covebasic: %d\n", len(ninoxBasicIncluded))
	fmt.Printf("from ninox, covebasic exlusion: %d\n", len(ninoxBasicExcluded))

	// generate indices for all clinialtrials entries and all ninox entries
	ninoxIndex := make(map[string]*ninox.Record)
	for i, record := range ninoxScreeningClinicaltrials {
		id := record.Field("nct_id")
		ninoxIndex[id] = &ninoxScreeningClinicaltrials[i]
	}

	// index of items in covebasic or covebasicexclude
	ninoxCoveBasicIndex := make(map[string]string)
	for _, record := range ninoxBasicIncluded {
		id := record.Field("source_id")
		source := record.Field("source")
		if source != "clinicaltrials.gov" {
			continue
		}
		ninoxCoveBasicIndex[id] = "included"
	}

	for _, record := range ninoxBasicExcluded {
		id := record.Field("source_id")
		source := record.Field("source")
		if source != "clinicaltrials.gov" {
			continue
		}
		ninoxCoveBasicIndex[id] = "excluded"
	}

	// iterate through all items in the source and compare it with the ninox data
	sourceIndex := make(map[string]*map[string]string)
	actionCounter := make(map[string]int)

	updates := make([][]string, 0)

	for i := range fromSource {

		record := fromSource[i]

		id := record["nct_id"]
		sourceIndex[id] = &fromSource[i]

		covebasic, ok := ninoxCoveBasicIndex[id]
		if !ok {
			covebasic = "not in covebasic"
		}

		// save current covebasic status
		record["covebasic"] = covebasic

		// skip all studies that are observational
		if strings.Contains(strings.ToLower(record["study_type"]), "observational") {

			// all studies that are not in ninox can be skipped
			_, ok := ninoxIndex[id]
			if !ok {
				// study should be skipped
				record["action"] = "skipped (observational)"
				actionCounter["skipped (observational)"]++

			} else {
				// study should be removed from ninox
				record["action"] = "remove from ninox (observational)"
				actionCounter["remove from ninox (observational)"]++
			}

			fromSource[i] = record
			continue
		}

		// check if a study is new
		ninoxRecord, ok := ninoxIndex[id]
		if !ok {
			record["action"] = "add to ninox"
			actionCounter["add to ninox"]++
			fromSource[i] = record
			continue
		}

		// compare all fields of the ninox record
		updateRequired := false
		for name, sourceValue := range record {
			if name == "action" || name == "covebasic" {
				continue
			}

			ninoxValue := ninoxRecord.Field(name)

			fullEq, partialEq := isEqual(sourceValue, ninoxValue)

			// nothing to do if the content is fully equivalent
			if fullEq {
				continue
			}

			// content is partially equal (lowercase and stripped out all non alphanum chars)
			if partialEq {
				continue
			}

			// otherwise an update is required
			updateRequired = true

			// add the update value to the updates list
			update := []string{
				"clinicaltrials.gov", id, name,
				ninoxValue, sourceValue,
			}
			updates = append(updates, update)

		}

		if updateRequired {
			record["action"] = "update in ninox"
			actionCounter["update in ninox"]++
			fromSource[i] = record
			continue
		}

		record["action"] = "equal to ninox record"
		actionCounter["equal to ninox record"]++
		fromSource[i] = record

	}

	for key, value := range actionCounter {
		fmt.Printf("Count: %03d, Action: %s\n", value, key)
	}

	// write the output to a new file
	file, err := os.Create(fmt.Sprintf("%s--records.csv", inputFile))
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer file.Close() // nolint:errcheck

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// write the header fields
	writer.Write(fields) // nolint:errcheck
	writer.Flush()

	// write out all the records
	for _, record := range fromSource {

		// create a new row
		row := make([]string, len(fields))

		// add all fields in the correct order (empty string if not set)
		for i, name := range fields {
			row[i] = record[name]
		}
		writer.Write(row) // nolint:errcheck
	}

	// nolint:errcheck
	writer.Flush()
	file.Close()

	// write the output to a new file
	file, err = os.Create(fmt.Sprintf("%s--updates.csv", inputFile))
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer file.Close() // nolint:errcheck

	writer = csv.NewWriter(file)
	writer.Comma = ';'

	// write the header fields
	writer.Write([]string{"source", "source_id", "field", "value_old", "value_new"}) // nolint:errcheck
	writer.WriteAll(updates)

	// nolint:errcheck
	writer.Flush()
	file.Close()

	return nil

}
