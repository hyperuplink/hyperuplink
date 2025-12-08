package forum

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Forum struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`

	CategoryID uuid.UUID `json:"category_id"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}
