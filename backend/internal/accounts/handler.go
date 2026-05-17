package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /accounts", h.listAccounts)
	mux.HandleFunc("POST /accounts", h.createAccount)
	mux.HandleFunc("GET /accounts/{id}", h.getAccount)
	mux.HandleFunc("PUT /accounts/{id}", h.updateAccount)
	mux.HandleFunc("DELETE /accounts/{id}", h.deactivateAccount)
}

type accountRequest struct {
	Name                string `json:"name"`
	Type                string `json:"type"`
	CurrencyCode        string `json:"currency_code"`
	InitialBalanceCents int64  `json:"initial_balance_cents"`
	IsActive            *bool  `json:"is_active"`
}

type accountResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	CurrencyCode        string    `json:"currency_code"`
	InitialBalanceCents int64     `json:"initial_balance_cents"`
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var req accountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if !isValidAccountType(req.Type) {
		writeJSONError(w, http.StatusBadRequest, "type must be one of: cash, checking, savings, credit_card")
		return
	}

	if strings.TrimSpace(req.CurrencyCode) == "" {
		writeJSONError(w, http.StatusBadRequest, "currency_code is required")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account, err := h.service.CreateAccount(r.Context(), CreateAccountInput{
		Name:                req.Name,
		Type:                req.Type,
		CurrencyCode:        strings.ToUpper(req.CurrencyCode),
		InitialBalanceCents: req.InitialBalanceCents,
		IsActive:            isActive,
	})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAccountResponse(*account))
}

func (h *Handler) getAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAccountResponse(*account))
}

func (h *Handler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListActiveAccounts(r.Context())
	if err != nil {
		writeMappedError(w, err)
		return
	}

	items := make([]accountResponse, 0, len(accounts))
	for _, account := range accounts {
		items = append(items, toAccountResponse(account))
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var req accountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if !isValidAccountType(req.Type) {
		writeJSONError(w, http.StatusBadRequest, "type must be one of: cash, checking, savings, credit_card")
		return
	}

	if strings.TrimSpace(req.CurrencyCode) == "" {
		writeJSONError(w, http.StatusBadRequest, "currency_code is required")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	account, err := h.service.UpdateAccount(r.Context(), UpdateAccountInput{
		ID:                  id,
		Name:                req.Name,
		Type:                req.Type,
		CurrencyCode:        strings.ToUpper(req.CurrencyCode),
		InitialBalanceCents: req.InitialBalanceCents,
		IsActive:            isActive,
	})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAccountResponse(*account))
}

func (h *Handler) deactivateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if err := h.service.DeactivateAccount(r.Context(), id); err != nil {
		writeMappedError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toAccountResponse(account repo.Account) accountResponse {
	return accountResponse{
		ID:                  uuid.UUID(account.ID.Bytes).String(),
		Name:                account.Name,
		Type:                account.Type,
		CurrencyCode:        account.CurrencyCode,
		InitialBalanceCents: account.InitialBalanceCents,
		IsActive:            account.IsActive,
		CreatedAt:           account.CreatedAt.Time,
		UpdatedAt:           account.UpdatedAt.Time,
	}
}

func parseUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue(key))
}

func isValidAccountType(accountType string) bool {
	switch accountType {
	case "cash", "checking", "savings", "credit_card":
		return true
	default:
		return false
	}
}

func writeMappedError(w http.ResponseWriter, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSONError(w, http.StatusNotFound, "resource not found")
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			writeJSONError(w, http.StatusConflict, "resource already exists")
			return
		case "23503":
			writeJSONError(w, http.StatusBadRequest, "invalid relation reference")
			return
		case "23514":
			writeJSONError(w, http.StatusBadRequest, "invalid value for constrained field")
			return
		}
	}

	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
