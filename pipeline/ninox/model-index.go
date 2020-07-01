package ninox

import "strings"

// Index contains information on covebasic inclusions/exclusions
type Index map[string]RecordInfo

type RecordInfo struct {
	ID     int
	Table  string
	Record *Record
}

// Set will save the given value in the index
func (i Index) Set(sourceID string, id int, table string, record *Record) {
	info := RecordInfo{
		ID:     id,
		Table:  table,
		Record: record,
	}

	sourceID = strings.ToLower(sourceID)
	sourceID = strings.TrimSpace(sourceID)
	i[sourceID] = info
}

// Get will check the index for the given value
func (i Index) Get(sourceID string) (RecordInfo, bool) {
	sourceID = strings.ToLower(sourceID)
	sourceID = strings.TrimSpace(sourceID)
	info, ok := i[sourceID]
	return info, ok
}
