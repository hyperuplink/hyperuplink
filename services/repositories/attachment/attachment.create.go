package attachment

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
)

func (repo *Repository) Create(model *attachment.Attachment) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO attachments (
		author_id
		,filename
		,mime_type
		,checksum
		,on_upload_hook_output
		,created_at
	) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,NOW()
	) RETURNING id`,
		model.AuthorID,
		model.Filename,
		model.MimeType,
		model.Checksum,
		model.OnUploadHookOutput,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
