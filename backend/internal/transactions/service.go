package transactions

import (
	"context"
	"strings"
	"time"

	"github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/dberrors"
	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/apperrors"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (*repo.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*repo.Transaction, error)
	ListTransactions(ctx context.Context, limit int32, offset int32) ([]repo.Transaction, error)
	ListTransactionsByAccount(ctx context.Context, input ListTransactionsByAccountInput) ([]repo.Transaction, error)
	UpdateTransaction(ctx context.Context, input UpdateTransactionInput) (*repo.Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
}

type CreateTransactionInput struct {
	AccountID       uuid.UUID
	CategoryID      *uuid.UUID
	PayeeName       string
	AmountCents     int64
	TransactionDate time.Time
	IsReconciled    bool
	Note            string
}

type UpdateTransactionInput struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	CategoryID      *uuid.UUID
	PayeeName       string
	AmountCents     int64
	TransactionDate time.Time
	IsReconciled    bool
	Note            string
}

type ListTransactionsByAccountInput struct {
	AccountID uuid.UUID
	FromDate  *time.Time
	ToDate    *time.Time
	Limit     int32
	Offset    int32
}

type svc struct {
	pool    *pgxpool.Pool
	queries *repo.Queries
}

func NewService(pool *pgxpool.Pool, queries *repo.Queries) Service {
	return &svc{pool: pool, queries: queries}
}

func (s *svc) CreateTransaction(ctx context.Context, input CreateTransactionInput) (*repo.Transaction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, dberrors.Map(err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	payeeID, err := resolvePayeeID(ctx, qtx, input.PayeeName)
	if err != nil {
		return nil, err
	}

	categoryID, err := resolveCategoryID(ctx, qtx, input.CategoryID)
	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateTransaction(ctx, repo.CreateTransactionParams{
		ID:              pgutil.ToPgUUID(uuid.New()),
		AccountID:       pgutil.ToPgUUID(input.AccountID),
		CategoryID:      categoryID,
		PayeeID:         payeeID,
		AmountCents:     input.AmountCents,
		TransactionDate: pgutil.Date(input.TransactionDate),
		IsReconciled:    input.IsReconciled,
		Note:            pgutil.NullableText(strings.TrimSpace(input.Note)),
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, dberrors.Map(err)
	}

	return &transaction, nil
}

func (s *svc) GetTransaction(ctx context.Context, id uuid.UUID) (*repo.Transaction, error) {
	transaction, err := s.queries.GetTransactionByID(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &transaction, nil
}

func (s *svc) ListTransactions(ctx context.Context, limit int32, offset int32) ([]repo.Transaction, error) {
	items, err := s.queries.ListTransactions(ctx, repo.ListTransactionsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return items, nil
}

func (s *svc) ListTransactionsByAccount(ctx context.Context, input ListTransactionsByAccountInput) ([]repo.Transaction, error) {
	if input.FromDate != nil && input.ToDate != nil {
		items, err := s.queries.ListTransactionsByAccountAndDateRange(ctx, repo.ListTransactionsByAccountAndDateRangeParams{
			AccountID:         pgutil.ToPgUUID(input.AccountID),
			TransactionDate:   pgutil.Date(*input.FromDate),
			TransactionDate_2: pgutil.Date(*input.ToDate),
			Limit:             input.Limit,
			Offset:            input.Offset,
		})
		if err != nil {
			return nil, dberrors.Map(err)
		}

		return items, nil
	}

	items, err := s.queries.ListTransactionsByAccount(ctx, repo.ListTransactionsByAccountParams{
		AccountID: pgutil.ToPgUUID(input.AccountID),
		Limit:     input.Limit,
		Offset:    input.Offset,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return items, nil
}

func (s *svc) UpdateTransaction(ctx context.Context, input UpdateTransactionInput) (*repo.Transaction, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, dberrors.Map(err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	payeeID, err := resolvePayeeID(ctx, qtx, input.PayeeName)
	if err != nil {
		return nil, err
	}

	categoryID, err := resolveCategoryID(ctx, qtx, input.CategoryID)
	if err != nil {
		return nil, err
	}

	transaction, err := qtx.UpdateTransaction(ctx, repo.UpdateTransactionParams{
		ID:              pgutil.ToPgUUID(input.ID),
		AccountID:       pgutil.ToPgUUID(input.AccountID),
		CategoryID:      categoryID,
		PayeeID:         payeeID,
		AmountCents:     input.AmountCents,
		TransactionDate: pgutil.Date(input.TransactionDate),
		IsReconciled:    input.IsReconciled,
		Note:            pgutil.NullableText(strings.TrimSpace(input.Note)),
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, dberrors.Map(err)
	}

	return &transaction, nil
}

func (s *svc) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	_, err := s.queries.DeleteTransaction(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return dberrors.Map(err)
	}

	return nil
}

func resolvePayeeID(ctx context.Context, qtx *repo.Queries, payeeName string) (pgtype.UUID, error) {
	trimmed := strings.TrimSpace(payeeName)
	if trimmed == "" {
		return pgtype.UUID{}, nil
	}
	trimmed = strings.ToLower(trimmed)

	payee, err := qtx.UpsertPayeeByNameCaseInsensitive(ctx, repo.UpsertPayeeByNameCaseInsensitiveParams{
		ID:   pgutil.ToPgUUID(uuid.New()),
		Name: trimmed,
	})
	if err != nil {
		return pgtype.UUID{}, dberrors.Map(err)
	}

	return payee.ID, nil
}

func resolveCategoryID(ctx context.Context, qtx *repo.Queries, categoryID *uuid.UUID) (pgtype.UUID, error) {
	if categoryID == nil {
		return pgtype.UUID{}, nil
	}

	resolved := pgutil.ToPgUUID(*categoryID)
	_, err := qtx.GetCategoryByID(ctx, resolved)
	if err == nil {
		return resolved, nil
	}

	mapped := dberrors.Map(err)
	if mapped == apperrors.ErrNotFound {
		return pgtype.UUID{}, apperrors.ErrInvalidReference
	}

	return pgtype.UUID{}, mapped
}
