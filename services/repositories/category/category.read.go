package category

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]category.Category, err error) {
	var mod []category.Category

	mod, err = database.Rows[category.Category](context.Background(), repo.db, qo.Query(
		`SELECT * FROM categories`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))

	return &mod, err
}

func (repo *Repository) GetBySlug(
	slug string,
	qo common.QueryOptions,
) (model *category.Category, err error) {
	return database.Row[category.Category](
		context.Background(),
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
	id uuid.UUID,
	qo common.QueryOptions,
) (model *category.Category, err error) {
	return database.Row[category.Category](
		context.Background(),
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
