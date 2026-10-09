package group

import (
	"context"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Delete(model *group.Group) (err error) {
	_, err = repo.db.Exec(context.Background(), `DELETE FROM groups WHERE id = $1`,
		model.ID,
	)
	return repo.db.ConvertError(err)
}
