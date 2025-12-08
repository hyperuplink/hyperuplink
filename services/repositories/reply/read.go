package reply

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/reply"
)

func (repo *Repository) AllForTopicUUID(
	id uuid.UUID,
) (model *[]reply.Reply, err error) {
	var rows pgx.Rows
	var mod []reply.Reply

	rows, err = repo.db.Query(`SELECT * FROM replies
		WHERE topic_id = $1
		AND deleted_at IS NULL
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[reply.Reply])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForTopicID(
	id string,
) (model *[]reply.Reply, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.AllForTopicUUID(uuID)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
) (model *reply.Reply, err error) {
	var rows pgx.Rows
	var mod reply.Reply

	rows, err = repo.db.Query(`SELECT * FROM replies
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[reply.Reply])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
) (model *reply.Reply, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID)
}
