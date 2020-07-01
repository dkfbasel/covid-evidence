package ninox

import (
	"dkfbasel.ch/covid-evidence/helpers"
)

// Field will return the value of the field as string if the field exists
// otherwise an empty string
func (r *Record) Field(name string) string {

	// return updated value if present
	updatedValue, ok := r.UpdatedFields[name]
	if ok {
		return helpers.AsString(updatedValue)
	}

	value, ok := r.Fields[name]
	if !ok {
		return ""
	}
	return helpers.AsString(value)
}
