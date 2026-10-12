package forum

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

func (repo *Repository) VAll(qo common.QueryOptions) (model *[]vforum.VForum, err error) {
	var mod []vforum.VForum

	mod, err = database.Rows[vforum.VForum](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vforums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))

	return &mod, err
}

func (repo *Repository) VAllForCategoryUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var mod []vforum.VForum

	mod, err = database.Rows[vforum.VForum](
		context.Background(),
		repo.db,
		qo.Query(`SELECT * FROM vforums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)

	return &mod, err
}

func (repo *Repository) VAllForCategoryID(
	id string,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.VAllForCategoryUUID(uuID, qo)
}

func (repo *Repository) VGetBySlug(
	slug string,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	return database.Row[vforum.VForum](
		context.Background(),
		repo.db,
		qo.Query(
			`SELECT * FROM vforums
			WHERE slug = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		slug,
	)
}

func (repo *Repository) VGetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	return database.Row[vforum.VForum](
		context.Background(),
		repo.db,
		qo.Query(
			`SELECT * FROM vforums
			WHERE id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
}
