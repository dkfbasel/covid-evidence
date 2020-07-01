package ninox

// Record is used to handle data from ninoxDB
type Record struct {
	BasicID       int                    `json:"-"`
	ID            int                    `json:"id,omitempty"`
	Sequence      int                    `json:"sequence,omitempty"`
	CreatedAt     string                 `json:"createdAt,omitempty"`
	CreatedBy     string                 `json:"createdBy,omitempty"`
	ModifiedAt    string                 `json:"modifiedAt,omitempty"`
	ModifiedBy    string                 `json:"modifiedBy,omitempty"`
	Fields        map[string]interface{} `json:"fields"`
	IsUpdated     bool                   `json:"-"`
	UpdatedFields map[string]interface{} `json:"-"`
}
