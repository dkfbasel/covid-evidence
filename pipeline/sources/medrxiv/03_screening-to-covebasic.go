package main

import (
	"fmt"

	"dkfbasel.ch/covid-evidence/helpers"
	"dkfbasel.ch/covid-evidence/ninox"
)

// convertRecords will convert the ictrp records to covebasic
func convertRecords(screeningRecords []ninox.Record, basicRecords []ninox.Record,
	basicIndex ninox.Index) []*ninox.Record {

	sourceName := "medRxiv"

	// initialize the updates/inserts
	updates := []*ninox.Record{}

	for _, s := range screeningRecords {

		sourceID := s.Field("ID")

		// skip all records that exist in ninox already
		info, ok := basicIndex[sourceID]
		if ok {
			fmt.Printf("record exists already: %s, %s\n", sourceID, info.Table)
			continue
		}

		// initialize a new record
		r := ninox.Record{}
		r.Fields = make(map[string]interface{})

		r.Fields["source"] = sourceName
		r.Fields["source_id"] = sourceID

		r.Fields["review_status"] = "prefilled automatically"

		r.Update("entry_type", "preprint", nil)

		r.Update("url", s.Field("rel_link"), nil)

		r.Update("title", s.Field("rel_title"), nil)
		r.Update("abstract", s.Field("rel_abs"), nil)
		r.Update("authors", s.Field("rel_authors"), nil)

		r.Update("doi", s.Field("rel_doi"), nil)

		r.Update("status_date", s.Field("rel_date"), helpers.ToIsoDate)

		// nothing to do, if the record was not changed
		if r.IsUpdated == false {
			continue
		}

		updates = append(updates, &r)
	}

	return updates

}
