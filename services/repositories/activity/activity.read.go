package activity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) MarkersForActor(
	kind activity.Kind,
	actorID uuid.UUID,
	subjectIDs []uuid.UUID,
	qo common.QueryOptions,
) (markers map[uuid.UUID]time.Time, err error) {
	markers = make(map[uuid.UUID]time.Time)

	if len(subjectIDs) == 0 {
		return markers, nil
	}

	var rows pgx.Rows
	rows, err = repo.db.Query(qo.Query(
		`SELECT
			 subject_id
			,updated_at
		FROM activities
		WHERE kind = $1 AND actor_id = $2 AND subject_id = ANY($3)`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		kind,
		actorID,
		subjectIDs,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var subjectID uuid.UUID
		var updatedAt pgtype.Timestamp

		if err = rows.Scan(&subjectID, &updatedAt); err != nil {
			return nil, repo.db.ConvertError(err)
		}

		if updatedAt.Valid {
			markers[subjectID] = updatedAt.Time
		}
	}

	return markers, repo.db.ConvertError(rows.Err())
}
