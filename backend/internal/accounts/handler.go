package accounts

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
	mux.HandleFunc("GET /accounts", h.listAccounts)
	mux.HandleFunc("POST /accounts", h.createAccount)
	mux.HandleFunc("GET /accounts/{id}/summary", h.getAccountSummary)
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

type accountSummaryResponse struct {
	AccountID              string `json:"account_id"`
	AccountName            string `json:"account_name"`
	BalanceCents           int64  `json:"balance_cents"`
	ReconciledBalanceCents int64  `json:"reconciled_balance_cents"`
	DifferenceCents        int64  `json:"difference_cents"`
}

func (h *Handler) createAccount(w http.ResponseWriter, r *http.Request) {
	var req accountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if !isValidAccountType(req.Type) {
		httputil.WriteJSONError(w, http.StatusBadRequest, "type must be one of: cash, checking, savings, credit_card")
		return
	}

	if strings.TrimSpace(req.CurrencyCode) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "currency_code is required")
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

	httputil.WriteJSON(w, http.StatusCreated, toAccountResponse(*account))
}

func (h *Handler) getAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toAccountResponse(*account))
}

func (h *Handler) getAccountSummary(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	summary, err := h.service.GetAccountSummary(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, accountSummaryResponse{
		AccountID:              summary.AccountID.String(),
		AccountName:            summary.AccountName,
		BalanceCents:           summary.BalanceCents,
		ReconciledBalanceCents: summary.ReconciledBalanceCents,
		DifferenceCents:        summary.DifferenceCents,
	})
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

	httputil.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	var req accountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	if !isValidAccountType(req.Type) {
		httputil.WriteJSONError(w, http.StatusBadRequest, "type must be one of: cash, checking, savings, credit_card")
		return
	}

	if strings.TrimSpace(req.CurrencyCode) == "" {
		httputil.WriteJSONError(w, http.StatusBadRequest, "currency_code is required")
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

	httputil.WriteJSON(w, http.StatusOK, toAccountResponse(*account))
}

func (h *Handler) deactivateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid account id")
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
