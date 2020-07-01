package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type coveBasic struct {
	CoveID              int    `json:"cove_id"`
	Source              string `json:"source"`
	ReviewStatus        string `json:"review_status"`
	ResultsAvailable    string `json:"results_available"`
	IpdSharing          string `json:"ipd_sharing"`
	InterventionType    string `json:"intervention_type"`
	InterventionName    string `json:"intervention_name"`
	NumberEnrollment    int    `json:"n_enrollment"`
	Country             string `json:"country"`
	Status              string `json:"status"`
	Randomized          string `json:"randomized"`
	NumberArms          int    `json:"n_arms"`
	Blinding            string `json:"blinding"`
	PopulationCondition string `json:"population_condition"`
	Control             string `json:"control"`
	OutPrimaryMeasure   string `json:"out_primary_measure"`
	StartDate           string `json:"start_date"`
	EndDate             string `json:"end_date"`
	SourceID            string `json:"source_id"`
	Title               string `json:"title"`
	URL                 string `json:"url"`
	IsCovid             string `json:"is_covid"`
	IsTrial             string `json:"is_trial"`
	IsObservational     string `json:"is_observational"`
	IsDuplicate         bool   `json:"is_duplicate,-"`
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("need path to preprocess")
	}

	fileName := os.Args[1]

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("could not read file content: %+v", err)
	}

	var dta []coveBasic
	err = json.Unmarshal(content, &dta)
	if err != nil {
		log.Fatalf("could not parse data: %+v", err)
	}

	var filtered []*coveBasic

	covidFiltered := 0
	trialFiltered := 0
	duplicateFiltered := 0
	sourceFiltered := 0
	idFiltered := 0

	// filter out all data that is not a trial
	for i, item := range dta {

		// skip all non covid items
		if item.IsCovid == "no" {
			covidFiltered++
			continue
		}

		// skip all non trial items
		if item.IsTrial == "no" {
			trialFiltered++
			continue
		}

		// skip all duplicates
		if item.IsDuplicate {
			duplicateFiltered++
			continue
		}

		if item.Source == "" {
			sourceFiltered++
			continue
		}

		// skip all items without covid id
		if item.CoveID == 0 {
			idFiltered++
			continue
		}

		filtered = append(filtered, &dta[i])
	}

	output, err := json.Marshal(&filtered)
	if err != nil {
		log.Fatalf("could not generate output: %+v", err)
	}

	outputName := strings.Replace(fileName, ".json", "_filtered.json", 1)
	err = ioutil.WriteFile(outputName, output, 0644)
	if err != nil {
		log.Fatalf("could not generate output file: %+v", err)
	}

	log.Printf("Input: %d, Filtered: %d", len(dta), len(filtered))
	log.Printf("Covid:   %d", covidFiltered)
	log.Printf("Trial:   %d", trialFiltered)
	log.Printf("Dupl.:   %d", duplicateFiltered)
	log.Printf("Source:  %d", sourceFiltered)
	log.Printf("ID:      %d", idFiltered)

}
