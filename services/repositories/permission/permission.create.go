package permission

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) Create(model *permission.Permission) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO permissions (
		 id
		,role
		,unit
		,forum_id
		,bits`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,$3
		,$4
		,$5`+
		// ,$6
		// ,$7
		// ,$8
		`) RETURNING id`,
		model.ID,
		model.Role,
		model.Unit,
		model.ForumID,
		model.Bits,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
