package categories

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/apperrors"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/httputil"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /categories", h.listCategories)
	mux.HandleFunc("POST /categories", h.createCategory)
	mux.HandleFunc("GET /categories/{id}", h.getCategory)
	mux.HandleFunc("PUT /categories/{id}", h.updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", h.deleteCategory)
}

type categoryRequest struct {
	Name string `json:"name"`
}

type categoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	category, err := h.service.CreateCategory(r.Context(), CreateCategoryInput{Name: req.Name})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, toCategoryResponse(*category))
}

func (h *Handler) getCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toCategoryResponse(*category))
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		writeMappedError(w, err)
		return
	}

	items := make([]categoryResponse, 0, len(categories))
	for _, c := range categories {
		items = append(items, toCategoryResponse(c))
	}

	httputil.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	category, err := h.service.UpdateCategoryName(r.Context(), UpdateCategoryNameInput{ID: id, Name: req.Name})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toCategoryResponse(*category))
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		writeMappedError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toCategoryResponse(category repo.Category) categoryResponse {
	return categoryResponse{
		ID:        uuid.UUID(category.ID.Bytes).String(),
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}
}

func parseUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue(key))
}

func writeMappedError(w http.ResponseWriter, err error) {
	if errors.Is(err, apperrors.ErrNotFound) {
		httputil.WriteJSONError(w, http.StatusNotFound, "resource not found")
		return
	}
	if errors.Is(err, apperrors.ErrAlreadyExists) {
		httputil.WriteJSONError(w, http.StatusConflict, "resource already exists")
		return
	}
	if errors.Is(err, apperrors.ErrInvalidReference) {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid relation reference")
		return
	}
	if errors.Is(err, apperrors.ErrConstraintViolation) {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid value for constrained field")
		return
	}

	httputil.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
}
