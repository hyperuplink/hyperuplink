package forum

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Create(model *forum.Forum) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO forums (
		 name
		,slug
		,position
		,category_id
		,description
		,created_at
		,updated_at
	) VALUES (
		 $1
		,$2
		,(SELECT COALESCE(MAX(position), 0) + 1 FROM forums WHERE category_id = $3 AND deleted_at IS NULL)
		,$3
		,$4
		,NOW()
		,NOW()
	) RETURNING id`,
		model.Name,
		model.Slug,
		model.CategoryID,
		model.Description,
	).Scan(&rowID)
	if err != nil {
		// TODO: Does it make sense to run this as a transaction?
		_, err = repo.db.Exec(`REFRESH MATERIALIZED VIEW vforums`)
	}
	return rowID, repo.db.ConvertError(err)
}
