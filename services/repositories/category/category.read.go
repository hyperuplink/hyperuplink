package category

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/memo"
)

func (repo *Repository) All(
	ctx context.Context,
	qo common.QueryOptions,
) (model []category.Category, err error) {
	query := qo.Query(
		`SELECT * FROM categories`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	)

	cached, err := memo.Get(repo.memo, "categories:"+query, func() ([]category.Category, error) {
		return database.Rows[category.Category](ctx, repo.db, query)
	})
	if err != nil {
		return nil, err
	}

	return slices.Clone(cached), nil
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
