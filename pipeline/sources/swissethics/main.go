package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"dkfbasel.ch/covid-evidence/ninox"
)

func main() {

	// fetch all items to be included in covebasic from ninox
	screeningRecords, err := ninox.FetchRecords(ninox.SwissethicsURL,
		`{"fields":{"cove_screening":"include"}}`)

	if err != nil {
		log.Printf("could not fetch screening records from ninox: %+v", err)
		return
	}

	// fetch all items from covebasic and index by medrxiv
	covebasicRecords, covebasicIndex, err := ninox.FetchCoveBasic("Ethics committees (CH)")

	log.Printf("fetched %d screening records", len(screeningRecords))
	log.Printf("fetched %d basic records", len(covebasicRecords))

	// convert the export to our basic table
	changes := convertRecords(
		screeningRecords,
		covebasicRecords,
		covebasicIndex,
	)

	log.Printf("have updates for %d records", len(changes))

	payload, _ := json.MarshalIndent(changes, "", "\t")
	ioutil.WriteFile("output.json", payload, 0777)

	// import the new basic records into ninox
	ninox.UpdateRecords(ninox.CoveBasicURL, changes)

}
