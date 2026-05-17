package payees

import (
	"context"

	"github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/dberrors"
	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/pgutil"
	"github.com/google/uuid"
)

type Service interface {
	CreatePayee(ctx context.Context, input CreatePayeeInput) (*repo.Payee, error)
	GetPayee(ctx context.Context, id uuid.UUID) (*repo.Payee, error)
	ListPayees(ctx context.Context) ([]repo.Payee, error)
	UpdatePayeeName(ctx context.Context, input UpdatePayeeNameInput) (*repo.Payee, error)
	DeletePayee(ctx context.Context, id uuid.UUID) error
}

type CreatePayeeInput struct {
	Name string
}

type UpdatePayeeNameInput struct {
	ID   uuid.UUID
	Name string
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreatePayee(ctx context.Context, input CreatePayeeInput) (*repo.Payee, error) {
	payee, err := s.repo.CreatePayee(ctx, repo.CreatePayeeParams{
		ID:   pgutil.ToPgUUID(uuid.New()),
		Name: input.Name,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &payee, nil
}

func (s *svc) GetPayee(ctx context.Context, id uuid.UUID) (*repo.Payee, error) {
	payee, err := s.repo.GetPayeeByID(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &payee, nil
}

func (s *svc) ListPayees(ctx context.Context) ([]repo.Payee, error) {
	payees, err := s.repo.ListPayees(ctx)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return payees, nil
}

func (s *svc) UpdatePayeeName(ctx context.Context, input UpdatePayeeNameInput) (*repo.Payee, error) {
	payee, err := s.repo.UpdatePayeeName(ctx, repo.UpdatePayeeNameParams{
		ID:   pgutil.ToPgUUID(input.ID),
		Name: input.Name,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &payee, nil
}

func (s *svc) DeletePayee(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.DeletePayee(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return dberrors.Map(err)
	}

	return nil
}
