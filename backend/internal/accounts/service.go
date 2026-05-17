package accounts

import (
	"context"
	"errors"

	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) (*repo.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*repo.Account, error)
	ListActiveAccounts(ctx context.Context) ([]repo.Account, error)
	UpdateAccount(ctx context.Context, input UpdateAccountInput) (*repo.Account, error)
	DeactivateAccount(ctx context.Context, id uuid.UUID) error
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
		ID:                  toPgUUID(uuid.New()),
		Name:                input.Name,
		Type:                input.Type,
		CurrencyCode:        input.CurrencyCode,
		InitialBalanceCents: input.InitialBalanceCents,
		IsActive:            input.IsActive,
	}

	account, err := s.repo.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *svc) GetAccount(ctx context.Context, id uuid.UUID) (*repo.Account, error) {
	account, err := s.repo.GetAccountByID(ctx, toPgUUID(id))
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *svc) ListActiveAccounts(ctx context.Context) ([]repo.Account, error) {
	return s.repo.ListActiveAccounts(ctx)
}

func (s *svc) UpdateAccount(ctx context.Context, input UpdateAccountInput) (*repo.Account, error) {
	params := repo.UpdateAccountParams{
		ID:                  toPgUUID(input.ID),
		Name:                input.Name,
		Type:                input.Type,
		CurrencyCode:        input.CurrencyCode,
		InitialBalanceCents: input.InitialBalanceCents,
		IsActive:            input.IsActive,
	}

	account, err := s.repo.UpdateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *svc) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeactivateAccount(ctx, toPgUUID(id))
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func fromPgUUID(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, errors.New("invalid pg UUID")
	}

	return uuid.UUID(id.Bytes), nil
}
