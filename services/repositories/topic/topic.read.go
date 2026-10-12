package topic

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	return database.Rows[topic.Topic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM topics`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
	)
}

func (repo *Repository) AllForForumUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	return database.Rows[topic.Topic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM topics WHERE forum_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllForForumID(
	id string,
	qo common.QueryOptions,
) (model []topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForForumUUID(uuID, qo)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	return database.Row[topic.Topic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM topics WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}

func (repo *Repository) GetByShortID(
	shortID string,
	qo common.QueryOptions,
) (model *topic.Topic, err error) {
	return database.Row[topic.Topic](context.Background(), repo.db, qo.Query(
		`SELECT * FROM topics WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}
