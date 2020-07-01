package main

import (
	"fmt"
	"log"

	"dkfbasel.ch/covid-evidence/helpers"
	"dkfbasel.ch/covid-evidence/ninox"
)

func main() {

	log.Println("fetching records")

	// fetch all records in the exclusion table
	records, err := ninox.FetchRecords(ninox.CoveBasicURL, "")
	if err != nil {
		fmt.Printf("could not fetch covebasic records: %+v\n", err)
		return
	}

	ctgov, err := ninox.FetchRecords(ninox.ClinicaltrialsURL, "")
	if err != nil {
		fmt.Printf("could not fetch ctgov records: %+v\n", err)
		return
	}

	fmt.Printf("got %d recors from ctgov\n", len(ctgov))

	ctgovIndex := make(map[string]*ninox.Record)

	for i, r := range ctgov {
		id := r.Field("nct_id")
		ctgovIndex[id] = &ctgov[i]
	}

	updated := []*ninox.Record{}

	for _, r := range records {

		newRecord := ninox.Record{}

		// initialize a new record
		newRecord.ID = r.ID
		newRecord.Fields = make(map[string]interface{})

		// date fields to check
		fields := []string{"status_date", "start_date", "end_date", "results_expected_date"}
		ctgovFields := []string{"date_last_update_posted", "date_started", "date_completed", ""}

		updatesAvailable := false

		for i, field := range fields {
			current := r.Field(field)

			// find corresponding ctgov entry
			source := r.Field("source")
			if source != "clinicaltrials.gov" {
				continue
			}

			nctId := r.Field("source_id")

			ctRecord, ok := ctgovIndex[nctId]
			if !ok {
				continue
			}

			if ctgovFields[i] == "" {
				continue
			}

			ctgovCurrent := ctRecord.Field(ctgovFields[i])
			transformed, _ := helpers.ToIsoDate(ctgovCurrent)

			if current != transformed {

				// skip items that are now empty
				if current != "" && transformed == "" {
					continue
				}

				updatesAvailable = true
				newRecord.Fields[field] = transformed

				if transformed == "" {
					transformed = "[EMPTY]"
				}

				// return

				fmt.Printf("% 6d: % 16s: % 18s: % 16s: % 16s\n", r.ID, nctId, field, current, transformed)
			}
		}

		if updatesAvailable {
			updated = append(updated, &newRecord)
		}

	}

	fmt.Printf("updates for %d records\n", len(updated))

	// var action string
	// fmt.Printf("Perform operation [n]: ")
	// fmt.Scanln(&action)

	// if action == "y" || action == "yes" {
	// 	fmt.Printf("change %d records\n", len(updated))
	// 	// update the given records
	// 	// ninox.UpdateRecords(ninox.CoveBasicURL, updated)
	// }

}
