package categories

import (
	"context"

	"github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/dberrors"
	repo "github.com/cubancodepath/zerobudget/backend/internal/adapters/postgresql/sqlc"
	"github.com/cubancodepath/zerobudget/backend/internal/shared/pgutil"
	"github.com/google/uuid"
)

type Service interface {
	CreateCategory(ctx context.Context, input CreateCategoryInput) (*repo.Category, error)
	GetCategory(ctx context.Context, id uuid.UUID) (*repo.Category, error)
	ListCategories(ctx context.Context) ([]repo.Category, error)
	UpdateCategoryName(ctx context.Context, input UpdateCategoryNameInput) (*repo.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type CreateCategoryInput struct {
	Name string
}

type UpdateCategoryNameInput struct {
	ID   uuid.UUID
	Name string
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateCategory(ctx context.Context, input CreateCategoryInput) (*repo.Category, error) {
	category, err := s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		ID:   pgutil.ToPgUUID(uuid.New()),
		Name: input.Name,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &category, nil
}

func (s *svc) GetCategory(ctx context.Context, id uuid.UUID) (*repo.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &category, nil
}

func (s *svc) ListCategories(ctx context.Context) ([]repo.Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return categories, nil
}

func (s *svc) UpdateCategoryName(ctx context.Context, input UpdateCategoryNameInput) (*repo.Category, error) {
	category, err := s.repo.UpdateCategoryName(ctx, repo.UpdateCategoryNameParams{
		ID:   pgutil.ToPgUUID(input.ID),
		Name: input.Name,
	})
	if err != nil {
		return nil, dberrors.Map(err)
	}

	return &category, nil
}

func (s *svc) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.DeleteCategory(ctx, pgutil.ToPgUUID(id))
	if err != nil {
		return dberrors.Map(err)
	}

	return nil
}
