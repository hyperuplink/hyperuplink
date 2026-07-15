package postevent

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *postevent.PostEvent, err error) {
	var rows pgx.Rows
	var mod postevent.PostEvent

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM postevents WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[postevent.PostEvent])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *postevent.PostEvent, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}

func (repo *Repository) TallyForTopicUUID(
	evType postevent.PostEventType,
	id uuid.UUID,
	qo common.QueryOptions,
) (tally map[int]int, err error) {
	var rows pgx.Rows

	tally = make(map[int]int)

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0

	rows, err = repo.db.Query(fmt.Sprintf(
		`SELECT selection, COUNT(*) AS votes FROM (%s) selections GROUP BY selection`,
		qoc.Query(
			`SELECT selection FROM postevents
			WHERE type = $1 AND target = $2 AND topic_id = $3`,
			common.QueryCapabilities{
				HasDeleted: true,
			}),
	),
		evType,
		postevent.Topic,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var selection int
		var votes int64

		if err = rows.Scan(&selection, &votes); err != nil {
			return nil, repo.db.ConvertError(err)
		}

		tally[selection] = int(votes)
	}

	return tally, repo.db.ConvertError(rows.Err())
}

func (repo *Repository) SelectionForTopicUUID(
	evType postevent.PostEventType,
	id uuid.UUID,
	authorID uuid.UUID,
	qo common.QueryOptions,
) (selection int, ok bool, err error) {
	err = repo.db.QueryRow(qo.Query(
		`SELECT selection FROM postevents
		WHERE type = $1 AND target = $2 AND topic_id = $3 AND author_id = $4`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		evType,
		postevent.Topic,
		id,
		authorID,
	).Scan(&selection)
	if err != nil {
		if err = repo.db.ConvertError(err); errors.Is(err, errs.ErrNoRows) {
			return -1, false, nil
		}
		return -1, false, err
	}

	return selection, true, nil
}
