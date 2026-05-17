package accounts

import (
	"context"

	"github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/dberrors"
	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/pgutil"
	"github.com/google/uuid"
)

type Service interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) (*repo.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*repo.Account, error)
	GetAccountSummary(ctx context.Context, id uuid.UUID) (*AccountSummary, error)
	ListActiveAccounts(ctx context.Context) ([]repo.Account, error)
	UpdateAccount(ctx context.Context, input UpdateAccountInput) (*repo.Account, error)
	DeactivateAccount(ctx context.Context, id uuid.UUID) error
}

type AccountSummary struct {
	AccountID              uuid.UUID
	AccountName            string
	BalanceCents           int64
	ReconciledBalanceCents int64
	DifferenceCents        int64
}

type CreateAccountInput struct {
	Name                string
	Type                string
	CurrencyCode        string
	InitialBalanceCents int64
	IsActive            bool
}

type UpdateAccountInput struct {
	ID                  uuid.UUID
	Name                string
	Type                string
	CurrencyCode        string
	InitialBalanceCents int64
	IsActive            bool
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateAccount(ctx context.Context, input CreateAccountInput) (*repo.Account, error) {
	params := repo.CreateAccountParams{
		ID:                  pgutil.ToPgUUID(uuid.New()),
		Name:                input.Name,
		Type:                input.Type,
		CurrencyCode:        input.CurrencyCode,
		InitialBalanceCents: input.InitialBalanceCents,
		IsActive:            input.IsActive,
	}

	account, err := s.repo.CreateAccount(ctx, params)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &account, nil
}

func (s *svc) GetAccount(ctx context.Context, id uuid.UUID) (*repo.Account, error) {
	account, err := s.repo.GetAccountByID(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &account, nil
}

func (s *svc) GetAccountSummary(ctx context.Context, id uuid.UUID) (*AccountSummary, error) {
	accountID := pgutil.ToPgUUID(id)

	balance, err := s.repo.GetAccountBalance(ctx, accountID)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	reconciled, err := s.repo.GetReconciledAccountBalance(ctx, accountID)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	summary := &AccountSummary{
		AccountID:              id,
		AccountName:            reconciled.AccountName,
		BalanceCents:           balance.BalanceCents,
		ReconciledBalanceCents: reconciled.BalanceCents,
	}
	summary.DifferenceCents = summary.BalanceCents - summary.ReconciledBalanceCents

	return summary, nil
}

func (s *svc) ListActiveAccounts(ctx context.Context) ([]repo.Account, error) {
	accounts, err := s.repo.ListActiveAccounts(ctx)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return accounts, nil
}

func (s *svc) UpdateAccount(ctx context.Context, input UpdateAccountInput) (*repo.Account, error) {
	params := repo.UpdateAccountParams{
		ID:                  pgutil.ToPgUUID(input.ID),
		Name:                input.Name,
		Type:                input.Type,
		CurrencyCode:        input.CurrencyCode,
		InitialBalanceCents: input.InitialBalanceCents,
		IsActive:            input.IsActive,
	}

	account, err := s.repo.UpdateAccount(ctx, params)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &account, nil
}

func (s *svc) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeactivateAccount(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return dberrors.Map(err)
	}

	return nil
}
