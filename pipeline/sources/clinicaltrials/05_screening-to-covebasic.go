package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"dkfbasel.ch/covid-evidence/helpers"
	"dkfbasel.ch/covid-evidence/ninox"
)

// ToCovebasic will transfer the records from the screening table to the covebasic table
func ToCovebasic() error {

	// fetch all items to be included in covebasic from ninox
	screeningRecords, err := ninox.FetchRecords(ninox.ClinicaltrialsURL,
		`{"fields":{"study_type":"Interventional"}}`)

	if err != nil {
		return fmt.Errorf("could not fetch screening records from ninox: %w", err)
	}

	// fetch all items from covebasic, indexed by clinicaltrials.gov
	covebasicRecords, covebasicIndex, err := ninox.FetchCoveBasic("clinicaltrials.gov")
	if err != nil {
		return fmt.Errorf("could not fetch covebasic records from ninox: %w", err)
	}

	log.Printf("fetched %d screening records", len(screeningRecords))
	log.Printf("fetched %d basic records", len(covebasicRecords))
	log.Printf("index contains %d records", len(covebasicIndex))

	// convert the export to our basic table
	changes := convertRecords(
		screeningRecords,
		covebasicRecords,
		covebasicIndex,
	)

	log.Printf("have updates for %d records", len(changes))

	var confirm string
	fmt.Print("Continue with update of covebasic [y/n]: ")
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "yes" {
		log.Println("abort")
		return nil
	}

	payload, _ := json.MarshalIndent(changes, "", "\t")
	ioutil.WriteFile("output.json", payload, 0777)

	// import the new basic records into ninox
	ninox.UpdateRecords(ninox.CoveBasicURL, changes)

	log.Println("update of covebasic completed")

	return nil

}

// convertRecords will convert the clinicaltrials.gov records to covebasic
func convertRecords(screeningRecords []ninox.Record, basicRecords []ninox.Record,
	basicIndex ninox.Index) []*ninox.Record {

	sourceName := "clinicaltrials.gov"

	// initialize the updates/inserts
	updates := []*ninox.Record{}

	for _, s := range screeningRecords {

		sourceID := s.Field("nct_id")

		// skip all records that exist in ninox already
		_, ok := basicIndex.Get(sourceID)
		if ok {
			// fmt.Printf("record exists already: %s, %s\n", sourceID, info.Table)
			continue
		}

		screening := s.Field("cove_screening")
		if screening == "1" || screening == "3" {
			continue
		}

		// initialize a new record
		r := ninox.Record{}
		r.Fields = make(map[string]interface{})

		r.Fields["source"] = sourceName
		r.Fields["source_id"] = sourceID

		r.Fields["review_status"] = "prefilled automatically"
		r.Fields["is_covid"] = "yes"
		r.Fields["is_trial"] = "yes"
		r.Fields["is_observational"] = "no"

		r.Update("entry_type", "registration", nil)

		r.Update("url", s.Field("nct_id"), func(value string) (interface{}, bool) {
			return fmt.Sprintf("https://clinicaltrials.gov/ct2/show/record/%s", value), true
		})

		r.Update("title", s.Fields["official_title"], nil)

		abstract := ""
		if !helpers.IsEmpty(s.Fields["brief_summary"]) {
			abstract = fmt.Sprintf("Brief summary:\n%s", helpers.AsString(s.Fields["brief_summary"]))
		}
		if !helpers.IsEmpty(s.Fields["detailed_description"]) {
			if abstract != "" {
				abstract = fmt.Sprintf("%s\n\n", abstract)
			}
			abstract = fmt.Sprintf("%s\n\nDetailed descriptions:\n%s", abstract,
				helpers.AsString(s.Fields["detailed_description"]))
		}

		abstract = strings.TrimSpace(abstract)

		r.Update("abstract", abstract, nil)

		r.Update("authors", "na", nil)
		r.Update("journal", "na", nil)
		r.Update("doi", "na", nil)

		r.Update("status", s.Fields["status"], helpers.ToLowerCase)

		r.Update("country", s.Fields["location_country"], func(country string) (interface{}, bool) {
			// country field may contain multiple countries separated by semicolon
			// -> use international if there are multiple countries
			// -> use the country name if it is the same multiple times
			if strings.Contains(country, ";") == false {
				return country, false
			}

			items := strings.Split(country, "; ")
			first := items[0]
			international := false
			for _, c := range items {
				if c != first {
					international = true
				}
			}
			if international {
				return "international", true
			}

			return first, true
		})

		// randomization
		r.Update("randomized", s.Fields["allocation"], helpers.ToLowerCase)

		// blinding
		r.Update("blinding", s.Fields["masking"], func(m string) (interface{}, bool) {
			if strings.HasPrefix(m, "None") {
				return "none", true
			}

			if strings.HasPrefix(m, "Double") || strings.HasPrefix(m, "Triple") || strings.HasPrefix(m, "Quadruple") {
				return "double blind", true
			}

			if strings.HasPrefix(m, "Single") {
				if strings.Contains(m, "Outcomes") {
					return "outcome only", true
				}
				return "single blind", true
			}

			return "", false
		})

		// longitudinal structure
		r.Update("longitudinal_structure", s.Fields["intervention_model"], helpers.ToLowerCase)

		// n_arms, calculated by the number of arm types
		r.Update("n_arms", s.Fields["arm_group_arm_group_type"], func(m string) (interface{}, bool) {
			count := strings.Count(m, "; ") + 1
			return count, true
		})

		// n_enrollment
		r.Update("n_enrollment", s.Fields["enrollment"], helpers.ToInt)

		// population_condition
		r.Update("population_condition", s.Fields["condition"], nil)

		// population_gender
		r.Update("population_gender", s.Fields["gender"], helpers.ToLowerCase)

		// skip population_age (difficult from min-max age)
		// r.Update("population_age", "TODO", nil)

		// skip intervention_type
		// r.Update("intervention_type"], "TODO", nil)

		// skip intervention and control (are in the same field)
		// r.Update("intervention_name"], "TODO", nil)
		// r.Update("control"], "TODO", nil)

		// out_primary_measure
		r.Update("out_primary_measure", s.Fields["primary_outcome_measure"], nil)

		// out_primary_desc
		r.Update("out_primary_desc", s.Fields["primary_outcome_description"], nil)

		// out_primary_timeframe
		r.Update("out_primary_timeframe", s.Fields["primary_outcome_time_frame"], nil)

		r.Update("start_date", s.Fields["date_started"], helpers.ToIsoDate)
		r.Update("end_date", s.Fields["date_completed"], helpers.ToIsoDate)

		// skip results_available
		// r.Update["results_available", "TODO", nil]

		// skip results_expected
		// r.Update["results_expected", "TODO", nil]

		// ipd_sharing
		r.Update("ipd_sharing", s.Fields["patient_data_sharing_ipd"], helpers.ToLowerCase)

		// publication
		r.Update("publication", s.Fields["publications_pmid"], nil)

		// out_secondary_measure
		r.Update("out_secondary_measure", s.Fields["secondary_outcome_measure"], nil)

		// out_secondary_desc
		r.Update("out_secondary_desc", s.Fields["secondary_outcome_description"], nil)

		// out_secondary_timeframe
		r.Update("out_secondary_timeframe", s.Fields["secondary_outcome_time_frame"], nil)

		// nothing to do, if the record was not changed
		if r.IsUpdated == false {
			continue
		}

		updates = append(updates, &r)
	}

	return updates

}
