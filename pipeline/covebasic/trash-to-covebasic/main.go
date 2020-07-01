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

	include := []*ninox.Record{}

	for i, record := range exclusionRecords {

		if record.Field("is_covid") == "yes" && record.Field("is_trial") == "yes" &&
			record.Field("is_duplicate") == "false" {
			include = append(include, &exclusionRecords[i])
		}

	}

	for i := range include {
		include[i].ID = 0
	}

	log.Printf("Items to move to covebasic table: %d", len(include))

	var action string
	fmt.Printf("Perform operation [n]: ")
	fmt.Scanln(&action)

	if action == "y" || action == "yes" {
		// import the exlusion records into the exlusion table
		ninox.UpdateRecords(ninox.CoveBasicURL, include)
	}

}
