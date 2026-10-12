package reply

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
)

func (repo *Repository) AllForTopicUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]reply.Reply, err error) {
	var mod []reply.Reply

	mod, err = database.Rows[reply.Reply](context.Background(), repo.db, qo.Query(
		`SELECT * FROM replies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)

	return &mod, err
}

func (repo *Repository) AllForTopicID(
	id string,
	qo common.QueryOptions,
) (model *[]reply.Reply, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.AllForTopicUUID(uuID, qo)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *reply.Reply, err error) {
	return database.Row[reply.Reply](context.Background(), repo.db, qo.Query(
		`SELECT * FROM replies WHERE id = $1`,
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
) (model *reply.Reply, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}

func (repo *Repository) GetByShortID(
	shortID string,
	qo common.QueryOptions,
) (model *reply.Reply, err error) {
	return database.Row[reply.Reply](context.Background(), repo.db, qo.Query(
		`SELECT * FROM replies WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}
