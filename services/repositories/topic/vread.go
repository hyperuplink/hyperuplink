package topic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/vtopic"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (repo *Repository) VAllForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vtopics WHERE forum_id = $1`,
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

func (repo *Repository) VAllForForumID(
	id string,
	qo common.QueryOptions,
) (model *[]vtopic.VTopic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.VAllForForumUUID(uuID, qo)
}

func (repo *Repository) VGetBySlug(
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
