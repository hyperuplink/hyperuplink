package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Update(ctx context.Context, model *group.Group) (err error) {
	_, err = database.Exec(ctx, repo.db, `UPDATE groups SET
		 name =                        $1
		,updated_at =                  NOW()
		WHERE id =                     $2`,
		model.Name,
		// WHERE id =
		model.ID,
	)
	return err
}
