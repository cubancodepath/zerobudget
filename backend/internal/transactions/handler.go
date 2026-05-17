package transactions

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/apperrors"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/httputil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const dateLayout = "2006-01-02"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /transactions", h.createTransaction)
	mux.HandleFunc("GET /transactions", h.listTransactions)
	mux.HandleFunc("GET /transactions/{id}", h.getTransaction)
	mux.HandleFunc("PUT /transactions/{id}", h.updateTransaction)
	mux.HandleFunc("DELETE /transactions/{id}", h.deleteTransaction)
	mux.HandleFunc("GET /accounts/{id}/transactions", h.listTransactionsByAccount)
}

type transactionRequest struct {
	AccountID       string  `json:"account_id"`
	CategoryID      *string `json:"category_id"`
	PayeeName       string  `json:"payee_name"`
	AmountCents     int64   `json:"amount_cents"`
	TransactionDate string  `json:"transaction_date"`
	IsReconciled    *bool   `json:"is_reconciled"`
	Note            string  `json:"note"`
}

type transactionResponse struct {
	ID              string     `json:"id"`
	AccountID       string     `json:"account_id"`
	CategoryID      *string    `json:"category_id"`
	PayeeID         *string    `json:"payee_id"`
	AmountCents     int64      `json:"amount_cents"`
	TransactionDate string     `json:"transaction_date"`
	IsReconciled    bool       `json:"is_reconciled"`
	Note            *string    `json:"note"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	input, err := parseTransactionInput(req)
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), input)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, toTransactionResponse(*transaction))
}

func (h *Handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	transaction, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toTransactionResponse(*transaction))
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.service.ListTransactions(r.Context(), limit, offset)
	if err != nil {
		writeMappedError(w, err)
		return
	}

	response := make([]transactionResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toTransactionResponse(item))
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) listTransactionsByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	limit, offset, err := parsePagination(r)
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	fromDate, err := parseOptionalDate(r.URL.Query().Get("from"))
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "from must use YYYY-MM-DD")
		return
	}

	toDate, err := parseOptionalDate(r.URL.Query().Get("to"))
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "to must use YYYY-MM-DD")
		return
	}

	if (fromDate == nil) != (toDate == nil) {
		httputil.WriteJSONError(w, http.StatusBadRequest, "from and to must be provided together")
		return
	}

	items, err := h.service.ListTransactionsByAccount(r.Context(), ListTransactionsByAccountInput{
		AccountID: accountID,
		FromDate:  fromDate,
		ToDate:    toDate,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	response := make([]transactionResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toTransactionResponse(item))
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	var req transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	input, err := parseTransactionInput(req)
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.service.UpdateTransaction(r.Context(), UpdateTransactionInput{
		ID:              id,
		AccountID:       input.AccountID,
		CategoryID:      input.CategoryID,
		PayeeName:       input.PayeeName,
		AmountCents:     input.AmountCents,
		TransactionDate: input.TransactionDate,
		IsReconciled:    input.IsReconciled,
		Note:            input.Note,
	})
	if err != nil {
		writeMappedError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toTransactionResponse(*updated))
}

func (h *Handler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDPathParam(r, "id")
	if err != nil {
		httputil.WriteJSONError(w, http.StatusBadRequest, "invalid transaction id")
		return
	}

	if err := h.service.DeleteTransaction(r.Context(), id); err != nil {
		writeMappedError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseTransactionInput(req transactionRequest) (CreateTransactionInput, error) {
	accountID, err := uuid.Parse(strings.TrimSpace(req.AccountID))
	if err != nil {
		return CreateTransactionInput{}, errors.New("account_id is required and must be a valid UUID")
	}

	if req.AmountCents == 0 {
		return CreateTransactionInput{}, errors.New("amount_cents is required and must not be zero")
	}

	transactionDate, err := time.Parse(dateLayout, strings.TrimSpace(req.TransactionDate))
	if err != nil {
		return CreateTransactionInput{}, errors.New("transaction_date is required and must use YYYY-MM-DD")
	}

	categoryID, err := parseOptionalUUID(req.CategoryID)
	if err != nil {
		return CreateTransactionInput{}, errors.New("category_id must be a valid UUID")
	}

	isReconciled := false
	if req.IsReconciled != nil {
		isReconciled = *req.IsReconciled
	}

	return CreateTransactionInput{
		AccountID:       accountID,
		CategoryID:      categoryID,
		PayeeName:       strings.TrimSpace(req.PayeeName),
		AmountCents:     req.AmountCents,
		TransactionDate: transactionDate,
		IsReconciled:    isReconciled,
		Note:            req.Note,
	}, nil
}

func parseOptionalUUID(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseOptionalDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := time.Parse(dateLayout, trimmed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parsePagination(r *http.Request) (int32, int32, error) {
	limit := int32(50)
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 {
			return 0, 0, errors.New("limit must be a positive integer")
		}
		limit = int32(parsedLimit)
	}

	offset := int32(0)
	if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			return 0, 0, errors.New("offset must be zero or a positive integer")
		}
		offset = int32(parsedOffset)
	}

	return limit, offset, nil
}

func toTransactionResponse(tx repo.Transaction) transactionResponse {
	return transactionResponse{
		ID:              uuid.UUID(tx.ID.Bytes).String(),
		AccountID:       uuid.UUID(tx.AccountID.Bytes).String(),
		CategoryID:      pgUUIDToStringPointer(tx.CategoryID),
		PayeeID:         pgUUIDToStringPointer(tx.PayeeID),
		AmountCents:     tx.AmountCents,
		TransactionDate: tx.TransactionDate.Time.Format(dateLayout),
		IsReconciled:    tx.IsReconciled,
		Note:            pgTextToStringPointer(tx.Note),
		CreatedAt:       tx.CreatedAt.Time,
		UpdatedAt:       tx.UpdatedAt.Time,
	}
}

func pgUUIDToStringPointer(value pgtype.UUID) *string {
	if !value.Valid {
		return nil
	}
	s := uuid.UUID(value.Bytes).String()
	return &s
}

func pgTextToStringPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func parseUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue(key))
}

func writeMappedError(w http.ResponseWriter, err error) {
	log.Printf("transactions handler error: %v", err)

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
