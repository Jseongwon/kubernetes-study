package entity

import (
	"time"
)

// JSONDocument represents a JSON document in our system
type JSONDocument struct {
	ID        string                 `json:"id"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// NewJSONDocument creates a new JSON document
func NewJSONDocument(id string, data map[string]interface{}) *JSONDocument {
	now := time.Now()
	return &JSONDocument{
		ID:        id,
		Data:      data,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the document data and timestamp
func (j *JSONDocument) Update(data map[string]interface{}) {
	j.Data = data
	j.UpdatedAt = time.Now()
}
