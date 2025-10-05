package repository

import (
	"context"
	"errors"
	"json-crud-service/internal/domain/entity"
)

var (
	ErrDocumentNotFound      = errors.New("document not found")
	ErrDocumentAlreadyExists = errors.New("document already exists")
)

// JSONRepository defines the interface for JSON document storage operations
type JSONRepository interface {
	// Create stores a new JSON document
	Create(ctx context.Context, doc *entity.JSONDocument) error
	
	// GetByID retrieves a JSON document by its ID
	GetByID(ctx context.Context, id string) (*entity.JSONDocument, error)
	
	// Update modifies an existing JSON document
	Update(ctx context.Context, doc *entity.JSONDocument) error
	
	// Delete removes a JSON document by its ID
	Delete(ctx context.Context, id string) error
	
	// List retrieves all JSON documents
	List(ctx context.Context) ([]*entity.JSONDocument, error)
}
