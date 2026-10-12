package topic

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

func (repo *Repository) All(
	ctx context.Context,
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	return database.Rows[topic.Topic](ctx, repo.db, qo.Query(
		`SELECT * FROM topics`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
	)
}

func (repo *Repository) AllForForumUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	return database.Rows[topic.Topic](ctx, repo.db, qo.Query(
		`SELECT * FROM topics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllForForumID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForForumUUID(ctx, uuID, qo)
}

func (repo *Repository) GetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	return database.Row[topic.Topic](ctx, repo.db, qo.Query(
		`SELECT * FROM topics WHERE id = $1`,
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
) (model *topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(ctx, uuID, qo)
}

func (repo *Repository) GetByShortID(
	ctx context.Context,
	shortID string,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	return database.Row[topic.Topic](ctx, repo.db, qo.Query(
		`SELECT * FROM topics WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}
