package attachment

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
)

func (repo *Repository) GetCategoryForAttachment(
	id uuid.UUID,
) (categoryID uuid.UUID, found bool, err error) {
	err = repo.db.QueryRow(`
		SELECT f.category_id FROM topics t
		JOIN forums f ON f.id = t.forum_id
		WHERE $1 = ANY(t.attachment_ids) AND t.deleted_at IS NULL
		UNION
		SELECT f.category_id FROM replies r
		JOIN topics t ON t.id = r.topic_id
		JOIN forums f ON f.id = t.forum_id
		WHERE $1 = ANY(r.attachment_ids)
		AND r.deleted_at IS NULL AND t.deleted_at IS NULL
		LIMIT 1`,
		id,
	).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return categoryID, false, nil
		}
		return categoryID, false, repo.db.ConvertError(err)
	}

	return categoryID, true, nil
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *attachment.Attachment, err error) {
	var rows pgx.Rows
	var mod attachment.Attachment

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM attachments WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[attachment.Attachment])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *attachment.Attachment, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}
