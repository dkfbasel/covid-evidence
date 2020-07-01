package main

import (
	"fmt"
	"log"

	"dkfbasel.ch/covid-evidence/ninox"
)

func main() {

	log.Println("fetching records")

	// fetch all records in the exclusion table
	exclusionRecords, err := ninox.FetchRecords(ninox.CoveBasicExlusionURL, "")
	if err != nil {
		fmt.Printf("could not fetch exclusion records: %+v\n", err)
		return
	}

	// create an index of the exclusion records using the source and sourceid as key
	exclusionIndex := make(map[string]int)
	for _, record := range exclusionRecords {
		key := record.Key()
		exclusionIndex[key] = record.ID
	}

	// free memory
	// nolint
	exclusionRecords = nil

	// fetch all records in the basic table
	basicRecords, err := ninox.FetchRecords(ninox.CoveBasicURL, "")
	if err != nil {
		fmt.Printf("could not fetch covebasic records: %+v\n", err)
		return
	}

	// find all items to exclude
	var exclude []*ninox.Record

	for i, record := range basicRecords {

		// skip all non covid items
		if record.Field("is_covid") == "no" {
			exclude = append(exclude, &basicRecords[i])
			continue
		}

		// skip all non trial items
		isTrial := record.Field("is_trial")
		isObservational := record.Field("is_observational")
		if isTrial == "no" {

			switch isObservational {
			case "no", "unclear", "yes":
				exclude = append(exclude, &basicRecords[i])
				continue
			}

		}

		// skip all duplicates
		if record.Field("is_duplicate") == "true" {
			exclude = append(exclude, &basicRecords[i])
			continue
		}
	}

	// initialize a list of records ids to be deleted from the basic table
	recordsToDelete := make([]int, len(exclude))

	// adapt the id for the exlusion table (i.e. updating existing or adding new)
	for i := range exclude {
		// save the basic id to delete the record later on
		recordsToDelete[i] = exclude[i].ID

		// match the key with existing items in the exlusion table and
		// overwrite the id of the exlusion records
		key := exclude[i].Key()
		id := exclusionIndex[key]
		exclude[i].ID = id
	}

	log.Printf("Items to move to exclusion table: %d", len(exclude))

	var action string
	fmt.Printf("Perform operation [n]: ")
	fmt.Scanln(&action)

	if action == "y" || action == "yes" {
		// import the exlusion records into the exlusion table
		ninox.UpdateRecords(ninox.CoveBasicExlusionURL, exclude)

		// delete the corresponding records from the basic table
		ninox.DeleteRecords(ninox.CoveBasicURL, recordsToDelete)
	}

}
