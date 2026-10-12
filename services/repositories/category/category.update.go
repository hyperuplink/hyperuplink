package category

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) Update(ctx context.Context, model *category.Category) (err error) {
	_, err = database.Exec(ctx, repo.db, `UPDATE categories SET
		 name =                        $1
		,slug =                        $2
		,updated_at =                  NOW()
		WHERE id =                     $3`,
		model.Name,
		model.Slug,
		// WHERE id =
		model.ID,
	)
	return err
}
