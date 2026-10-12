package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Delete(ctx context.Context, model *group.Group) (err error) {
	_, err = database.Exec(ctx, repo.db, `DELETE FROM groups WHERE id = $1`,
		model.ID,
	)
	return err
}
