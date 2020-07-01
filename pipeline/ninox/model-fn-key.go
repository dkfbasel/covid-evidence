package ninox

import (
	"fmt"
)

// Key will return the key of the item in ninox (combined from source and source_id)
func (r *Record) Key() string {
	source := r.Field("source")
	sourceID := r.Field("source_id")
	if source == "" && sourceID == "" {
		return ""
	}
	if sourceID == "" {
		return source
	}
	if source == "" {
		return sourceID
	}

	key := fmt.Sprintf("%s::%s", source, sourceID)
	return key
}
