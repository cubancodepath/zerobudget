package dberrors

import (
	"errors"

	"github.com/cubancodepath/zerobudget/backend/internal/shared/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Map(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperrors.ErrAlreadyExists
		case "23503":
			return apperrors.ErrInvalidReference
		case "23514":
			return apperrors.ErrConstraintViolation
		}
	}

	return err
}
