package forum

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) All(qo common.QueryOptions) (model *[]forum.Forum, err error) {
	var mod []forum.Forum

	mod, err = database.Rows[forum.Forum](context.Background(), repo.db, qo.Query(
		`SELECT * FROM forums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))

	return &mod, err
}

func (repo *Repository) AllForCategoryUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]forum.Forum, err error) {
	var mod []forum.Forum

	mod, err = database.Rows[forum.Forum](
		context.Background(),
		repo.db,
		qo.Query(`SELECT * FROM forums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)

	return &mod, err
}

func (repo *Repository) AllForCategoryID(
	id string,
	qo common.QueryOptions,
) (model *[]forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForCategoryUUID(uuID, qo)
}

func (repo *Repository) GetBySlug(
	slug string,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	return database.Row[forum.Forum](
		context.Background(),
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
	id uuid.UUID,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	return database.Row[forum.Forum](
		context.Background(),
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
	id string,
	qo common.QueryOptions,
) (model *forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}
