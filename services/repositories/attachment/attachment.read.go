package attachment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
)

func (repo *Repository) GetCategoryForAttachment(
	ctx context.Context,
	id uuid.UUID,
) (categoryID uuid.UUID, found bool, err error) {
	categoryID, err = database.Value[uuid.UUID](ctx, repo.db, `
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
	)
	if err != nil {
		if errors.Is(err, errs.ErrNoRows) {
			return categoryID, false, nil
		}
		return categoryID, false, err
	}

	return categoryID, true, nil
}

func (repo *Repository) GetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *attachment.Attachment, err error) {
	return database.Row[attachment.Attachment](ctx, repo.db, qo.Query(
		`SELECT * FROM attachments WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) GetByID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model *attachment.Attachment, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(ctx, uuID, qo)
}

func (repo *Repository) AllByUUIDs(
	ctx context.Context,
	ids []uuid.UUID,
	qo common.QueryOptions,
) (model []attachment.Attachment, err error) {
	return database.Rows[attachment.Attachment](ctx, repo.db, qo.Query(
		`SELECT * FROM attachments WHERE id = ANY($1)`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		ids,
	)
}
