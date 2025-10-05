package handler

import (
	"encoding/json"
	"fmt"
	"json-crud-service/internal/usecase"
	"json-crud-service/pkg/response"
	"net/http"
)

// JSONHandler handles HTTP requests for JSON document operations
type JSONHandler struct {
	usecase *usecase.JSONUsecase
}

// NewJSONHandler creates a new JSON handler instance
func NewJSONHandler(usecase *usecase.JSONUsecase) *JSONHandler {
	return &JSONHandler{
		usecase: usecase,
	}
}

// CreateDocumentRequest represents the request payload for creating a document
type CreateDocumentRequest struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// UpdateDocumentRequest represents the request payload for updating a document
type UpdateDocumentRequest struct {
	Data map[string]interface{} `json:"data"`
}

// CreateDocument handles POST /documents
func (h *JSONHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Errorf("invalid JSON payload: %v", err))
		return
	}

	doc, err := h.usecase.CreateDocument(r.Context(), req.ID, req.Data)
	if err != nil {
		if err == usecase.ErrInvalidID || err == usecase.ErrInvalidData {
			response.BadRequest(w, err)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.Created(w, doc)
}

// GetDocument handles GET /documents/{id}
func (h *JSONHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	doc, err := h.usecase.GetDocument(r.Context(), id)
	if err != nil {
		if err == usecase.ErrInvalidID {
			response.BadRequest(w, err)
		} else if err == usecase.ErrDocumentNotFound {
			response.NotFound(w, err)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.OK(w, doc)
}

// UpdateDocument handles PUT /documents/{id}
func (h *JSONHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Errorf("invalid JSON payload: %v", err))
		return
	}

	doc, err := h.usecase.UpdateDocument(r.Context(), id, req.Data)
	if err != nil {
		if err == usecase.ErrInvalidID || err == usecase.ErrInvalidData {
			response.BadRequest(w, err)
		} else if err == usecase.ErrDocumentNotFound {
			response.NotFound(w, err)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.OK(w, doc)
}

// DeleteDocument handles DELETE /documents/{id}
func (h *JSONHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.usecase.DeleteDocument(r.Context(), id)
	if err != nil {
		if err == usecase.ErrInvalidID {
			response.BadRequest(w, err)
		} else if err == usecase.ErrDocumentNotFound {
			response.NotFound(w, err)
		} else {
			response.InternalServerError(w, err)
		}
		return
	}

	response.NoContent(w)
}

// ListDocuments handles GET /documents
func (h *JSONHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := h.usecase.ListDocuments(r.Context())
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	response.OK(w, docs)
}

// HealthCheck handles GET /health
func (h *JSONHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.OK(w, map[string]string{
		"status":  "healthy",
		"service": "json-crud-service",
	})
}

// SetupRoutes sets up all the HTTP routes
func (h *JSONHandler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("POST /documents", h.CreateDocument)
	mux.HandleFunc("GET /documents", h.ListDocuments)
	mux.HandleFunc("GET /documents/{id}", h.GetDocument)
	mux.HandleFunc("PUT /documents/{id}", h.UpdateDocument)
	mux.HandleFunc("DELETE /documents/{id}", h.DeleteDocument)
}
