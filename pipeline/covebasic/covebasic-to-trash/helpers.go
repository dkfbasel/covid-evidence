package main

import "fmt"

// return string representation of value
func asString(value interface{}) string {

	if value == nil {
		return ""
	}

	switch t := value.(type) {
	case int, int32, int64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%.0f", t)
	case bool:
		return fmt.Sprintf("%t", t)
	default:
		return fmt.Sprintf("%s", t)
	}

}

func getKey(record *NinoxRecord) string {
	source := asString(record.Fields["source"])
	sourceID := asString(record.Fields["source_id"])
	key := fmt.Sprintf("%s::%s", source, sourceID)
	return key
}
