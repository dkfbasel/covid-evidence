package main

import (
	"fmt"
	"log"
	"strings"

	"dkfbasel.ch/covid-evidence/ninox"
)

func main() {

	// now fetch all data from ninox
	screeningRecords, err := ninox.FetchRecords(ninox.IctrpURL, "")
	if err != nil {
		log.Printf("could not fetch screening records from ninox: %+v", err)
		return
	}

	// fetch all items from covebasic and index the ictrp sources
	covebasicRecords, covebasicIndex, err := ninox.FetchCoveBasic("ictrp")

	log.Printf("fetched %d screening records", len(screeningRecords))
	log.Printf("fetched %d basic records", len(covebasicRecords))

	removeFromCoveBasic := []int{}
	updateScreening := []*ninox.Record{}

	// find all items in the icrtrp table that have study type observational and
	// are listed in the covebasic table and have review status "prefilled automatically"
	for _, r := range screeningRecords {

		studyType := strings.ToLower(r.Field("Study type"))

		// nothing to do if observational is not in the study type
		if strings.Contains(studyType, "observational") == false {
			continue
		}

		sourceID := r.Field("TrialID")

		// nothing to do if the source is not in clinicaltrials
		info, ok := covebasicIndex.Get(sourceID)
		if !ok {
			continue
		}

		// skip all items that are in the exclusions table
		if info.Table == ninox.CoveBasicExlusionsTable {
			continue
		}

		reviewStatus := info.Record.Field("review_status")
		if reviewStatus != "prefilled automatically" {
			fmt.Printf("%s: %05d, %s\n", sourceID, info.ID, reviewStatus)
			continue
		}

		fmt.Printf("modify: %s\n", sourceID)

		// remove the record from cove basic
		removeFromCoveBasic = append(removeFromCoveBasic, info.ID)

		// adapt the screening record status to exlude
		sRecord := ninox.Record{}
		sRecord.ID = r.ID
		sRecord.Fields = make(map[string]interface{})
		sRecord.Fields["cove_screening"] = "automatic exclusion"
		updateScreening = append(updateScreening, &sRecord)

	}

	log.Printf("update for screening records:  %03d", len(updateScreening))
	log.Printf("delete from covebasic records: %03d", len(removeFromCoveBasic))

	// import the new basic records into ninox
	// if len(updateScreening) > 0 {
	// 	ninox.UpdateRecords(ninox.IctrpURL, updateScreening)
	// }

	// if len(removeFromCoveBasic) > 0 {
	// 	ninox.DeleteRecords(ninox.CoveBasicURL, removeFromCoveBasic)
	// }

}
