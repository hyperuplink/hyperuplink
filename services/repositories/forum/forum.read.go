package forum

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) All(ctx context.Context, qo common.QueryOptions) (model []forum.Forum, err error) {
	return database.Rows[forum.Forum](ctx, repo.db, qo.Query(
		`SELECT * FROM forums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
}

func (repo *Repository) AllForCategoryUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []forum.Forum, err error) {
	return database.Rows[forum.Forum](
		ctx,
		repo.db,
		qo.Query(`SELECT * FROM forums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
}

func (repo *Repository) AllForCategoryID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model []forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForCategoryUUID(ctx, uuID, qo)
}

func (repo *Repository) GetBySlug(
	ctx context.Context,
	slug string,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	return database.Row[forum.Forum](
		ctx,
		repo.db,
		qo.Query(
			`SELECT * FROM forums
			WHERE slug = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		slug,
	)
}

func (repo *Repository) GetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	return database.Row[forum.Forum](
		ctx,
		repo.db,
		qo.Query(
			`SELECT * FROM forums
			WHERE id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
}

func (repo *Repository) GetByID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(ctx, uuID, qo)
}
