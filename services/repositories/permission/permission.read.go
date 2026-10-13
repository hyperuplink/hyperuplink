package permission

import (
	"context"
	"slices"

	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/memo"
)

func (repo *Repository) all(ctx context.Context) ([]permission.Permission, error) {
	return memo.Get(repo.memo, "permissions", func() ([]permission.Permission, error) {
		return database.Rows[permission.Permission](ctx, repo.db, `SELECT * FROM permissions
			WHERE deleted_at IS NULL`)
	})
}

func (repo *Repository) All(ctx context.Context) (model []permission.Permission, err error) {
	cached, err := repo.all(ctx)
	if err != nil {
		return nil, err
	}

	return slices.Clone(cached), nil
}

func (repo *Repository) AllForGroup(
	ctx context.Context,
	groupID string,
) (model []permission.Permission, err error) {
	return database.Rows[permission.Permission](ctx, repo.db, `SELECT * FROM permissions
		WHERE group_id = $1
		AND deleted_at IS NULL`,
		groupID,
	)
}

func (repo *Repository) GetDefault(ctx context.Context) (model *permission.Permission, err error) {
	return repo.GetFor(ctx, pgtype.Text{}, pgtype.UUID{})
}

func (repo *Repository) GetFor(
	ctx context.Context,
	groupID pgtype.Text,
	categoryID pgtype.UUID,
) (model *permission.Permission, err error) {
	return database.Row[permission.Permission](ctx, repo.db, `SELECT * FROM permissions
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2
		AND deleted_at IS NULL
		LIMIT 1`,
		groupID,
		categoryID,
	)
}
