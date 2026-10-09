package group

import (
	"context"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Create(model *group.Group) (rowID string, err error) {
	err = repo.db.QueryRow(context.Background(), `INSERT INTO groups (
		 id
		,name
		,created_at
		,updated_at
		,deleted_at
		) VALUES (
		 $1
		,$2
		,NOW()
		,NOW()
		,NULL
		) RETURNING id`,
		model.ID,
		model.Name,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
