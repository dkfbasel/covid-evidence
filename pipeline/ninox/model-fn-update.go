package ninox

import (
	"fmt"

	"dkfbasel.ch/covid-evidence/helpers"
)

func (r *Record) Update(fieldName string, sourceField interface{}, fn handlerFunc) {

	// define the name of the associated certainty field
	fieldNameCertainty := fmt.Sprintf("%s_certainty", fieldName)

	sourceValue := helpers.AsString(sourceField)
	currentValue := helpers.AsString(r.Fields[fieldName])
	currentCertainty := helpers.AsString(r.Fields[fieldNameCertainty])

	// is the content generated
	isGenerated := false

	fieldUpdated := false

	// use a custom handler function for the variable if specified
	if fn != nil {
		sourceField, isGenerated = fn(sourceValue)
		sourceValue = helpers.AsString(sourceField)
	}

	// empty fields
	if sourceValue == "" {
		// set certainty to prefilled if not set or set to generated
		// nothing to do if already set to prefilled
		if r.Fields[fieldNameCertainty] == "" || r.Fields[fieldNameCertainty] == "generated" {
			r.Fields[fieldNameCertainty] = "prefilled"
			r.IsUpdated = true
			fieldUpdated = true
		}
	}

	// check if the content has changed and update the value if necessary
	if sourceValue != currentValue {

		if currentCertainty == "human" {
			// log.Printf("extracted value has changed %s:\n  - old: %s\n  - new: %s", fieldName, currentValue, sourceValue)
			delete(r.Fields, fieldName)
			delete(r.Fields, fieldNameCertainty)
			return
		}

		if currentCertainty == "verified" {
			// log.Printf("verified value has changed %s:\n  - old: %s\n  - new: %s", fieldName, currentValue, sourceValue)
			delete(r.Fields, fieldName)
			delete(r.Fields, fieldNameCertainty)
			return
		}

		// log.Printf("field has changed %s:\n - %s\n - %s", fieldName, currentValue, sourceValue)

		r.Fields[fieldName] = sourceField
		r.IsUpdated = true
		fieldUpdated = true
		if isGenerated {
			r.Fields[fieldNameCertainty] = "generated"
		} else {
			r.Fields[fieldNameCertainty] = "prefilled"
		}
	}

	// remove the field information if it is not updated
	if !fieldUpdated {
		delete(r.Fields, fieldName)
		delete(r.Fields, fieldNameCertainty)
	}

}

// handlerFunc is used to handle specific fields
type handlerFunc func(value string) (interface{}, bool)
