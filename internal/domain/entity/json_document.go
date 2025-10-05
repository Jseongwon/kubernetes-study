package entity

import (
	"time"
)

// JSONDocument represents a JSON document in our system
type JSONDocument struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// NewJSONDocument creates a new JSON document
func NewJSONDocument(id, docType, version string, data map[string]interface{}) *JSONDocument {
	now := time.Now()
	return &JSONDocument{
		ID:        id,
		Type:      docType,
		Version:   version,
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

// UpdateWithVersion updates the document data, version and timestamp
func (j *JSONDocument) UpdateWithVersion(data map[string]interface{}, newVersion string) {
	j.Data = data
	j.Version = newVersion
	j.UpdatedAt = time.Now()
}
