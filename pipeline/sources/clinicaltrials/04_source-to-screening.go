package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"dkfbasel.ch/covid-evidence/ninox"
)

// Import will import new studies from clinicaltrials gov into the ninox database
func Import(inputFile string) error {

	inputFile = fmt.Sprintf("%s--records", inputFile)

	// compile a list of all fields that should be checked in ninox
	fields := []string{"action", "covebasic"}
	fieldInputRowIndex := make(map[string]int)

	fieldInputRowIndex["action"] = 0
	fieldInputRowIndex["covebasic"] = 1

	inputFieldIndex := 2
	for _, field := range fieldMap {
		if field.Ninox != "" {
			fields = append(fields, field.Ninox)
			fieldInputRowIndex[field.Ninox] = inputFieldIndex
			inputFieldIndex++
		}
	}

	// first load data from the local file
	fromSource, err := loadFromFile(inputFile, fieldInputRowIndex)
	if err != nil {
		return fmt.Errorf("could not load records from file: %+w", err)
	}

	// initialize records to import in ninox
	toNinox := []*ninox.Record{}

	// get the current date for imported and last_update datestamp
	currentDate := time.Now().Format("2006-01-02")

	// go through all rows of the csv file (records from ctgov)
	for _, row := range fromSource {

		// skip everything that should not be added to ninox screening
		if row["action"] != "add to ninox" {
			fmt.Printf("action: %s\n", row["action"])
			continue
		}

		// skip everything that is already in the covebasic table
		if row["covebasic"] != "not in covebasic" {
			fmt.Printf("covebasic: %s\n", row["covebasic"])
			continue
		}

		// create a new record
		record := ninox.Record{}

		// initialize fields
		record.Fields = make(map[string]interface{})

		record.Fields["cove_import_date"] = currentDate
		record.Fields["cove_update_date"] = currentDate

		// add all fields to the record
		for key, value := range row {
			if key == "action" || key == "covebasic" {
				continue
			}

			if key == "enrollment" {
				asInt, _ := strconv.Atoi(value)
				record.Fields[key] = asInt
				continue
			}

			record.Fields[key] = value
		}

		toNinox = append(toNinox, &record)
	}

	fmt.Printf("import records to ninox: %d\n", len(toNinox))

	var confirm string
	fmt.Print("Continue with import into screening [y/n]: ")
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "yes" {
		log.Println("abort")
		return nil
	}

	// update the records in ninox
	ninox.UpdateRecords(ninox.ClinicaltrialsURL, toNinox)

	return nil
}
