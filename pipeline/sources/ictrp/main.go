package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"dkfbasel.ch/covid-evidence/ninox"
)

func main() {

	// now ictrp records that should be included
	ninoxScreening, err := ninox.FetchRecords(ninox.IctrpURL,
		`{"fields":{"cove_screening":"include"}}`)

	if err != nil {
		log.Printf("could not fetch clnicaltrials records from ninox: %+v", err)
		return
	}

	ninoxBasic, ninoxBasicIndex, err := ninox.FetchCoveBasic("ICTRP", "clinicaltrials.gov")
	if err != nil {
		log.Printf("could not fetch covebasic records: %+v", err)
	}

	log.Printf("fetched %d screening records", len(ninoxScreening))
	log.Printf("fetched %d basic records", len(ninoxBasic))

	// convert the export to our basic table
	changes := convertRecords(
		ninoxScreening,
		ninoxBasicIndex,
	)

	log.Printf("have updates for %d records", len(changes))

	var confirm string
	fmt.Print("Continue with update of covebasic [y/n]: ")
	fmt.Scanln(&confirm) // nolint
	if confirm != "y" && confirm != "yes" {
		log.Println("abort")
		return
	}

	// save the data to an output file
	payload, _ := json.MarshalIndent(changes, "", "\t")
	ioutil.WriteFile("output.json", payload, 0777)

	// import the new basic records into ninox
	ninox.UpdateRecords(ninox.CoveBasicURL, changes)

}
