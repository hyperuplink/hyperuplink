package forum

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

func (repo *Repository) VAll(qo common.QueryOptions) (model *[]vforum.VForum, err error) {
	var rows pgx.Rows
	var mod []vforum.VForum

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vforums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForCategoryUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var rows pgx.Rows
	var mod []vforum.VForum

	rows, err = repo.db.Query(
		qo.Query(`SELECT * FROM vforums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForCategoryID(
	id string,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.VAllForCategoryUUID(uuID, qo)
}

func (repo *Repository) VGetBySlug(
	slug string,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	var rows pgx.Rows
	var mod vforum.VForum

	rows, err = repo.db.Query(
		qo.Query(
			`SELECT * FROM vforums
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

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VGetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vforum.VForum, err error) {
	var rows pgx.Rows
	var mod vforum.VForum

	rows, err = repo.db.Query(
		qo.Query(
			`SELECT * FROM vforums
			WHERE id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}
