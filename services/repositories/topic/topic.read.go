package topic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/topic"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]topic.Topic, err error) {
	var rows pgx.Rows
	var mod []topic.Topic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM topics`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[topic.Topic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]topic.Topic, err error) {
	var rows pgx.Rows
	var mod []topic.Topic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM topics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[topic.Topic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForForumID(
	id string,
	qo common.QueryOptions,
) (model *[]topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.AllForForumUUID(uuID, qo)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	var rows pgx.Rows
	var mod topic.Topic

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM topics WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[topic.Topic])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}
