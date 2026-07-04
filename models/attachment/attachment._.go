package attachment

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Attachment struct {
	ID uuid.UUID `json:"id"`

	AuthorID uuid.UUID `json:"author_id"`

	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Checksum string `json:"checksum"`

	OnUploadHookOutput string `json:"on_upload_hook_output"`

	CreatedAt   pgtype.Timestamp `json:"created_at"`
	ModeratedAt pgtype.Timestamp `json:"moderated_at"`
	SpammedAt   pgtype.Timestamp `json:"spammed_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}
