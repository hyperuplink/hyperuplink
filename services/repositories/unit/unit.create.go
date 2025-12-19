package unit

import (
	"github.com/mrusme/hyperuplink/models/unit"
)

func (repo *Repository) Create(model *unit.Unit) (rowID string, err error) {
	err = repo.db.QueryRow(`INSERT INTO units (
		 id
		,created_at
		,updated_at
		,deleted_at
		) VALUES (
		,$1
		,NOW()
		,NOW()
		,NULL
		) RETURNING id`,
		model.ID,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
