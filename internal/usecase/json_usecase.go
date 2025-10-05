package usecase

import (
	"context"
	"errors"
	"json-crud-service/internal/domain/entity"
	"json-crud-service/internal/domain/repository"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrInvalidID        = errors.New("invalid document ID")
	ErrInvalidData      = errors.New("invalid document data")
)

// JSONUsecase defines the business logic for JSON document operations
type JSONUsecase struct {
	repo repository.JSONRepository
}

// NewJSONUsecase creates a new JSON usecase instance
func NewJSONUsecase(repo repository.JSONRepository) *JSONUsecase {
	return &JSONUsecase{
		repo: repo,
	}
}

// CreateDocument creates a new JSON document
func (u *JSONUsecase) CreateDocument(ctx context.Context, id string, data map[string]interface{}) (*entity.JSONDocument, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	
	if data == nil {
		return nil, ErrInvalidData
	}
	
	doc := entity.NewJSONDocument(id, data)
	
	if err := u.repo.Create(ctx, doc); err != nil {
		return nil, err
	}
	
	return doc, nil
}

// GetDocument retrieves a JSON document by ID
func (u *JSONUsecase) GetDocument(ctx context.Context, id string) (*entity.JSONDocument, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	
	doc, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if doc == nil {
		return nil, ErrDocumentNotFound
	}
	
	return doc, nil
}

// UpdateDocument updates an existing JSON document
func (u *JSONUsecase) UpdateDocument(ctx context.Context, id string, data map[string]interface{}) (*entity.JSONDocument, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	
	if data == nil {
		return nil, ErrInvalidData
	}
	
	doc, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if doc == nil {
		return nil, ErrDocumentNotFound
	}
	
	doc.Update(data)
	
	if err := u.repo.Update(ctx, doc); err != nil {
		return nil, err
	}
	
	return doc, nil
}

// DeleteDocument removes a JSON document
func (u *JSONUsecase) DeleteDocument(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID
	}
	
	// Check if document exists
	doc, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if doc == nil {
		return ErrDocumentNotFound
	}
	
	return u.repo.Delete(ctx, id)
}

// ListDocuments retrieves all JSON documents
func (u *JSONUsecase) ListDocuments(ctx context.Context) ([]*entity.JSONDocument, error) {
	return u.repo.List(ctx)
}
