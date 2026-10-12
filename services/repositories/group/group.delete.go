package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Delete(model *group.Group) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `DELETE FROM groups WHERE id = $1`,
		model.ID,
	)
	return err
}
