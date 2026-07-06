package category

import (
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]category.Category, err error) {
	var rows pgx.Rows
	var mod []category.Category

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM categories`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[category.Category])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetBySlug(
	slug string,
	qo common.QueryOptions,
) (model *category.Category, err error) {
	var rows pgx.Rows
	var mod category.Category

	rows, err = repo.db.Query(
		qo.Query(
			`SELECT * FROM categories
			WHERE slug = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		slug,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])

	return &mod, repo.db.ConvertError(err)
}

// func (repo *Repository) GetByUUID(
// 	id uuid.UUID,
// ) (model *category.Category, err error) {
// 	var rows pgx.Rows
// 	var mod category.Category
//
// 	rows, err = repo.db.Query(`SELECT * FROM categories
// 		WHERE id = $1
// 		AND deleted_at IS NULL
// 		LIMIT 1`,
// 		id)
// 	if err != nil {
// 		return nil, repo.db.ConvertError(err)
// 	}
//
// 	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[category.Category])
//
// 	return &mod, repo.db.ConvertError(err)
// }
//
// func (repo *Repository) GetByID(
// 	id string,
// ) (model *category.Category, err error) {
// 	var uuID uuid.UUID
//
// 	if uuID, err = uuid.Parse(id); err != nil {
// 		return nil, repo.db.ConvertError(err)
// 	}
//
// 	return repo.GetByUUID(uuID)
// }
