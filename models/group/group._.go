package group

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}
