package activity

import (
	"fmt"
	"slices"
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/models/activity"
)

const (
	columnsPerRecord       int = 7
	maxRecordsPerStatement int = 1000
)

func (repo *Repository) values(recs []activity.Record) (vals string, args []any) {
	var groups []string

	for i, rec := range recs {
		base := i * columnsPerRecord
		var params []string
		for n := range columnsPerRecord {
			params = append(params, fmt.Sprintf("$%d", base+n+1))
		}
		groups = append(groups, fmt.Sprintf("(%s)", strings.Join(params, ",")))

		var dedupeKey any
		if rec.DedupeKey != "" {
			dedupeKey = rec.DedupeKey
		}

		count := rec.Count
		if count < 1 {
			count = 1
		}

		var context any
		if len(rec.Context) > 0 {
			context = rec.Context
		}

		args = append(args,
			rec.Kind,
			rec.ActorID,
			rec.Subject,
			rec.SubjectID,
			dedupeKey,
			count,
			context,
		)
	}

	return strings.Join(groups, ","), args
}

func (repo *Repository) Upsert(recs []activity.Record) (err error) {
	for chunk := range slices.Chunk(recs, maxRecordsPerStatement) {
		vals, args := repo.values(chunk)

		if _, err = repo.db.Exec(fmt.Sprintf(`INSERT INTO activities (
			 kind
			,actor_id
			,subject
			,subject_id
			,dedupe_key
			,count
			,context
		) VALUES %s
		ON CONFLICT (kind, actor_id, dedupe_key)
		WHERE dedupe_key IS NOT NULL AND deleted_at IS NULL
		DO UPDATE SET
			 count = activities.count + EXCLUDED.count
			,updated_at = NOW()`, vals),
			args...,
		); err != nil {
			return repo.db.ConvertError(err)
		}
	}

	return nil
}

func (repo *Repository) Append(recs []activity.Record) (err error) {
	for chunk := range slices.Chunk(recs, maxRecordsPerStatement) {
		vals, args := repo.values(chunk)

		if _, err = repo.db.Exec(fmt.Sprintf(`INSERT INTO activities (
			 kind
			,actor_id
			,subject
			,subject_id
			,dedupe_key
			,count
			,context
		) VALUES %s`, vals),
			args...,
		); err != nil {
			return repo.db.ConvertError(err)
		}
	}

	return nil
}
