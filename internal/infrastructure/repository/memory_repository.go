package repository

import (
	"context"
	"json-crud-service/internal/domain/entity"
	"json-crud-service/internal/domain/repository"
	"sync"
)

// MemoryJSONRepository implements JSONRepository interface using in-memory storage
type MemoryJSONRepository struct {
	documents map[string]*entity.JSONDocument
	mu        sync.RWMutex
}

// NewMemoryJSONRepository creates a new in-memory JSON repository
func NewMemoryJSONRepository() repository.JSONRepository {
	return &MemoryJSONRepository{
		documents: make(map[string]*entity.JSONDocument),
	}
}

// Create stores a new JSON document
func (r *MemoryJSONRepository) Create(ctx context.Context, doc *entity.JSONDocument) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if document already exists
	if _, exists := r.documents[doc.ID]; exists {
		return repository.ErrDocumentAlreadyExists
	}

	// Create a copy to avoid external modifications
	docCopy := &entity.JSONDocument{
		ID:        doc.ID,
		Type:      doc.Type,
		Version:   doc.Version,
		Data:      make(map[string]interface{}),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}

	// Deep copy the data map
	for k, v := range doc.Data {
		docCopy.Data[k] = v
	}

	r.documents[doc.ID] = docCopy
	return nil
}

// GetByID retrieves a JSON document by its ID
func (r *MemoryJSONRepository) GetByID(ctx context.Context, id string) (*entity.JSONDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doc, exists := r.documents[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid external modifications
	docCopy := &entity.JSONDocument{
		ID:        doc.ID,
		Type:      doc.Type,
		Version:   doc.Version,
		Data:      make(map[string]interface{}),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}

	// Deep copy the data map
	for k, v := range doc.Data {
		docCopy.Data[k] = v
	}

	return docCopy, nil
}

// Update modifies an existing JSON document
func (r *MemoryJSONRepository) Update(ctx context.Context, doc *entity.JSONDocument) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if document exists
	if _, exists := r.documents[doc.ID]; !exists {
		return repository.ErrDocumentNotFound
	}

	// Create a copy to avoid external modifications
	docCopy := &entity.JSONDocument{
		ID:        doc.ID,
		Type:      doc.Type,
		Version:   doc.Version,
		Data:      make(map[string]interface{}),
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}

	// Deep copy the data map
	for k, v := range doc.Data {
		docCopy.Data[k] = v
	}

	r.documents[doc.ID] = docCopy
	return nil
}

// Delete removes a JSON document by its ID
func (r *MemoryJSONRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if document exists
	if _, exists := r.documents[id]; !exists {
		return repository.ErrDocumentNotFound
	}

	delete(r.documents, id)
	return nil
}

// List retrieves all JSON documents
func (r *MemoryJSONRepository) List(ctx context.Context) ([]*entity.JSONDocument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	documents := make([]*entity.JSONDocument, 0, len(r.documents))
	for _, doc := range r.documents {
		// Create a copy for each document
		docCopy := &entity.JSONDocument{
			ID:        doc.ID,
			Type:      doc.Type,
			Version:   doc.Version,
			Data:      make(map[string]interface{}),
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		}

		// Deep copy the data map
		for k, v := range doc.Data {
			docCopy.Data[k] = v
		}

		documents = append(documents, docCopy)
	}

	return documents, nil
}
