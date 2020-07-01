package main

import (
	"fmt"

	"dkfbasel.ch/covid-evidence/helpers"
	"dkfbasel.ch/covid-evidence/ninox"
)

// convertRecords will convert the ictrp records to covebasic
func convertRecords(screeningRecords []ninox.Record, basicRecords []ninox.Record,
	basicIndex ninox.Index) []*ninox.Record {

	// initialize the updates/inserts
	updates := []*ninox.Record{}

	for _, s := range screeningRecords {

		sourceID := s.Field("Project ID")

		// skip all records that exist in ninox already
		info, ok := basicIndex[sourceID]
		if ok {
			fmt.Printf("record exists already: %s, %s\n", sourceID, info.Table)

			if info.Record.Field("review_status") != "prefilled automatically" {
				continue
			}

		}

		// initialize a new record
		r := ninox.Record{}
		r.Fields = make(map[string]interface{})

		if ok {
			r.ID = info.ID
		}

		r.Fields["source"] = "Ethics committees (CH)"
		r.Fields["source_id"] = sourceID

		r.Fields["review_status"] = "prefilled automatically"
		r.Fields["is_trial"] = "yes"

		r.Update("entry_type", "ethics", nil)

		r.Update("title", s.Field("Project Title"), nil)
		r.Update("authors", s.Field("Principal Investigator"), nil)

		r.Update("country", "Switzerland", nil)

		r.Update("status_date", s.Field("Date final decision"), helpers.ToIsoDate)

		r.Update("funding", s.Field("Sponsor"), nil)

		r.Fields["extraction_comment"] = s.Field("Type of Project")

		// nothing to do, if the record was not changed
		if r.IsUpdated == false {
			continue
		}

		updates = append(updates, &r)
	}

	return updates

}
