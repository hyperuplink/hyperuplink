package topic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

func (repo *Repository) VAllCountForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (total int64, err error) {
	var rows pgx.Rows

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	rows, err = repo.db.Query(qoc.Query(
		`SELECT COUNT(id) AS total FROM vtopics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	pag, err = pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}
	total = pag["total"].(int64)

	return total, nil
}

func (repo *Repository) VAllCountForReadableSlugs(
	slugs []string,
	qo common.QueryOptions,
) (total int64, err error) {
	var rows pgx.Rows

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0

	base := `SELECT COUNT(id) AS total FROM vtopics`
	var args []any
	if slugs != nil {
		base += ` WHERE category_slug = ANY($1)`
		args = append(args, slugs)
	}

	rows, err = repo.db.Query(qoc.Query(
		base,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		args...,
	)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	pag, err = pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}
	total = pag["total"].(int64)

	return total, nil
}

func (repo *Repository) VAllForReadableSlugs(
	slugs []string,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	total, err = repo.VAllCountForReadableSlugs(slugs, qo)
	if err != nil {
		return nil, 0, err
	}

	base := `SELECT * FROM vtopics`
	var args []any
	if slugs != nil {
		base += ` WHERE category_slug = ANY($1)`
		args = append(args, slugs)
	}

	rows, err = repo.db.Query(qo.Query(
		base,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		args...,
	)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, total, repo.db.ConvertError(err)
}

func (repo *Repository) VAll(
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	total, err = repo.VAllCountForForumUUID(id, qo)

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, total, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForForumID(
	id string,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	return repo.VAllForForumUUID(uuID, qo)
}

func (repo *Repository) VGetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod vtopic.VTopic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VGetByShortID(
	shortID string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod vtopic.VTopic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VGetByForumUUIDSlug(
	forumUUID uuid.UUID,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod vtopic.VTopic

	rows, err = repo.db.Query(
		qo.Query(
			`SELECT * FROM vtopics
			WHERE forum_id = $1 AND slug = $2`,
			common.QueryCapabilities{
				HasSpammed: true,
				HasDeleted: true,
			},
		),
		forumUUID,
		slug,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VGetByForumIDSlug(
	forumID string,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(forumID); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.VGetByForumUUIDSlug(uuID, slug, qo)
}

func (repo *Repository) VGetBySlugs(
	forumSlug string,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod vtopic.VTopic

	rows, err = repo.db.Query(
		qo.Query(
			`SELECT * FROM vtopics
			WHERE forum_slug = $1 AND slug = $2`,
			common.QueryCapabilities{
				HasSpammed: true,
				HasDeleted: true,
			},
		),
		forumSlug,
		slug,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForAuthorUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics WHERE author_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, repo.db.ConvertError(err)
}
