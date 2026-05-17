package payees

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
	mux.HandleFunc("GET /payees", h.listPayees)
	mux.HandleFunc("POST /payees", h.createPayee)
	mux.HandleFunc("GET /payees/{id}", h.getPayee)
	mux.HandleFunc("PUT /payees/{id}", h.updatePayee)
	mux.HandleFunc("DELETE /payees/{id}", h.deletePayee)
}

type payeeRequest struct {
	Name string `json:"name"`
}

type payeeResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) createPayee(w http.ResponseWriter, r *http.Request) {
	var req payeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	payee, err := h.service.CreatePayee(r.Context(), CreatePayeeInput{Name: req.Name})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, toPayeeResponse(*payee))
}

func (h *Handler) getPayee(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid payee id")
		return
	}

	payee, err := h.service.GetPayee(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toPayeeResponse(*payee))
}

func (h *Handler) listPayees(w http.ResponseWriter, r *http.Request) {
	payees, err := h.service.ListPayees(r.Context())
	if err != nil {
		writeMappedError(w, err)
		return
	}

	items := make([]payeeResponse, 0, len(payees))
	for _, p := range payees {
		items = append(items, toPayeeResponse(p))
	}

	httputil.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) updatePayee(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid payee id")
		return
	}

	var req payeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	payee, err := h.service.UpdatePayeeName(r.Context(), UpdatePayeeNameInput{ID: id, Name: req.Name})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toPayeeResponse(*payee))
}

func (h *Handler) deletePayee(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid payee id")
		return
	}

	if err := h.service.DeletePayee(r.Context(), id); err != nil {
		writeMappedError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toPayeeResponse(payee repo.Payee) payeeResponse {
	return payeeResponse{
		ID:        uuid.UUID(payee.ID.Bytes).String(),
		Name:      payee.Name,
		CreatedAt: payee.CreatedAt.Time,
		UpdatedAt: payee.UpdatedAt.Time,
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
