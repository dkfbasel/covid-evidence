package main

import (
	"fmt"
	"strings"

	"dkfbasel.ch/covid-evidence/ninox"
)

// convertRecords will convert the ictrp records to covebasic
func convertRecords(screeningRecords []ninox.Record, basicIndex ninox.Index) []*ninox.Record {

	const sourceName = "ICTRP"

	// initialize the updates/inserts
	updates := []*ninox.Record{}

	for _, s := range screeningRecords {

		sourceID := s.Field("TrialID")

		// skip all records that exist in ninox already
		_, ok := basicIndex.Get(sourceID)
		if ok {
			fmt.Printf("record exists already: %s, %s\n", sourceID)
			continue
		}

		// initialize a new record
		r := ninox.Record{}
		r.Fields = make(map[string]interface{})
		r.Fields["source"] = sourceName
		r.Fields["source_id"] = sourceID
		r.Fields["review_status"] = "prefilled automatically"

		r.Update("entry_type", "registration", nil)

		r.Update("url", s.Field("web address"), nil)

		r.Update("title", s.Field("Scientific title"), nil)

		r.Update("corresp_author_lastname", s.Field("Contact Lastname"), nil)
		r.Update("corresp_author_email", s.Field("Contact Email"), nil)

		r.Update("status", s.Field("Recruitment Status"), toLowerCase)
		r.Update("status_date", s.Field("Last Refreshed On"), toIsoDate)

		r.Update("country", s.Field("Countries"), func(country string) (interface{}, bool) {
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

		r.Update("randomized", s.Field("Study design"), toLowerCase)

		r.Update("population_condition", s.Field("condition"), nil)

		r.Update("intervention_name", s.Field("Intervention"), nil)

		r.Update("out_primary_measure", s.Field("Primary outcome"), nil)

		r.Update("start_date", s.Field("Date enrollement"), toIsoDate)

		// results_available if a results url is given
		r.Update("results_available", s.Field("results url link"), func(url string) (interface{}, bool) {
			if url == "" {
				return "no", true
			}
			return "yes", true
		})

		r.Update("inclusion_criteria", s.Field("Inclusion Criteria"), nil)
		r.Update("exclusion_criteria", s.Field("Exclusion Criteria"), nil)

		// nothing to do, if the record was not changed
		if r.IsUpdated == false {
			continue
		}

		updates = append(updates, &r)
	}

	return updates

}
