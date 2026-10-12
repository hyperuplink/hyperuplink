package forum

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

func (repo *Repository) VAll(ctx context.Context, qo common.QueryOptions) (model []vforum.VForum, err error) {
	return database.Rows[vforum.VForum](ctx, repo.db, qo.Query(
		`SELECT * FROM vforums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
}

func (repo *Repository) VAllForCategoryUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []vforum.VForum, err error) {
	return database.Rows[vforum.VForum](
		ctx,
		repo.db,
		qo.Query(`SELECT * FROM vforums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
}

func (repo *Repository) VAllForCategoryID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model []vforum.VForum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.VAllForCategoryUUID(ctx, uuID, qo)
}

func (repo *Repository) VGetBySlug(
	ctx context.Context,
	slug string,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	return database.Row[vforum.VForum](
		ctx,
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
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	return database.Row[vforum.VForum](
		ctx,
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
