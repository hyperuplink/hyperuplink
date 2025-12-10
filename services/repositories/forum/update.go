package forum

import (
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Update(model *forum.Forum) (err error) {
	_, err = repo.db.Exec(`UPDATE forums SET
		 name =                        $1
		,slug =                        $2
		,position = CASE
			WHEN category_id <> $3 OR category_id IS DISTINCT FROM $3 THEN (SELECT COALESCE(MAX(position), 0) + 1 FROM forums WHERE category_id = $3 AND deleted_at IS NULL)
			ELSE position
		END
		,category_id =                 $3
		,description =                 $4
		,updated_at =                 NOW()
		WHERE id =                     $5`,
		model.Name,
		model.Slug,
		model.CategoryID,
		model.Description,
		model.ID,
	)
	if err != nil {
		// TODO: Does it make sense to run this as a transaction?
		_, err = repo.db.Exec(`REFRESH MATERIALIZED VIEW vforums`)
	}
	return repo.db.ConvertError(err)
}
