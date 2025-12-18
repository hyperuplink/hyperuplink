package unit

import (
	"github.com/mrusme/hyperuplink/models/unit"
)

func (repo *Repository) Create(model *unit.Unit) (rowID string, err error) {
	err = repo.db.QueryRow(`INSERT INTO units (
		 id`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		,$1`+
		// ,$2
		// ,$3
		// ,$4
		`) RETURNING id`,
		model.ID,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
