package category

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) All(
	ctx context.Context,
	qo common.QueryOptions,
) (model []category.Category, err error) {
	return database.Rows[category.Category](ctx, repo.db, qo.Query(
		`SELECT * FROM categories`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
}

func (repo *Repository) GetBySlug(
	ctx context.Context,
	slug string,
	qo common.QueryOptions,
) (model *category.Category, err error) {
	return database.Row[category.Category](
		ctx,
		repo.db,
		qo.Query(
			`SELECT * FROM categories
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
) (model *category.Category, err error) {
	return database.Row[category.Category](
		ctx,
		repo.db,
		qo.Query(
			`SELECT * FROM categories
			WHERE id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
}
