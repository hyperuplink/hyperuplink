package topic

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

func (repo *Repository) VAllCountForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (total int64, err error) {

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	return database.Value[int64](context.Background(), repo.db, qoc.Query(
		`SELECT COUNT(id) AS total FROM vtopics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) VAllCountForReadableSlugs(
	slugs []string,
	qo common.QueryOptions,
) (total int64, err error) {
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

	return database.Value[int64](context.Background(), repo.db, qoc.Query(
		base,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		args...,
	)
}

func (repo *Repository) VAllForReadableSlugs(
	slugs []string,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
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

	mod, err = database.Rows[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		base,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		args...,
	)

	return &mod, total, err
}

func (repo *Repository) VAll(
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var mod []vtopic.VTopic

	mod, err = database.Rows[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vtopics`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
	)

	return &mod, err
}

func (repo *Repository) VAllForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
	var mod []vtopic.VTopic

	if total, err = repo.VAllCountForForumUUID(id, qo); err != nil {
		return nil, 0, err
	}

	mod, err = database.Rows[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vtopics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)

	return &mod, total, err
}

func (repo *Repository) VAllForForumID(
	id string,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, total int64, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, 0, database.ConvertError(err)
	}

	return repo.VAllForForumUUID(uuID, qo)
}

func (repo *Repository) VGetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	return database.Row[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vtopics WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) VGetByShortID(
	shortID string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	return database.Row[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vtopics WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}

func (repo *Repository) VGetByForumUUIDSlug(
	forumUUID uuid.UUID,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	return database.Row[vtopic.VTopic](
		context.Background(),
		repo.db,
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
}

func (repo *Repository) VGetByForumIDSlug(
	forumID string,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(forumID); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.VGetByForumUUIDSlug(uuID, slug, qo)
}

func (repo *Repository) VGetBySlugs(
	forumSlug string,
	slug string,
	qo common.QueryOptions,
) (model *vtopic.VTopic, err error) {
	return database.Row[vtopic.VTopic](
		context.Background(),
		repo.db,
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
}

func (repo *Repository) VAllForAuthorUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var mod []vtopic.VTopic

	mod, err = database.Rows[vtopic.VTopic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM vtopics WHERE author_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)

	return &mod, err
}
