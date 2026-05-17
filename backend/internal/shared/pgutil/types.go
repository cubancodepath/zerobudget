package pgutil

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NullableText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{String: value, Valid: true}
}

func Date(value time.Time) pgtype.Date {
	return pgtype.Date{Time: value, Valid: true}
}
