package unit

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Unit struct {
	ID string `json:"id"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}
