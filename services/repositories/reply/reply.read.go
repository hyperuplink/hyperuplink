package reply

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
)

func (repo *Repository) AllForTopicUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []reply.Reply, err error) {
	return database.Rows[reply.Reply](ctx, repo.db, qo.Query(
		`SELECT * FROM replies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllForTopicID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model []reply.Reply, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForTopicUUID(ctx, uuID, qo)
}

func (repo *Repository) GetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *reply.Reply, err error) {
	return database.Row[reply.Reply](ctx, repo.db, qo.Query(
		`SELECT * FROM replies WHERE id = $1`,
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
) (model *reply.Reply, err error) {
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
) (model *reply.Reply, err error) {
	return database.Row[reply.Reply](ctx, repo.db, qo.Query(
		`SELECT * FROM replies WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}
