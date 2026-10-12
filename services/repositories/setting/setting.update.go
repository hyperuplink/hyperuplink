package setting

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func Update[T any](
	ctx context.Context,
	repo *Repository,
	model *setting.Setting[T],
) (err error) {
	_, err = database.Exec(ctx, repo.db, `UPDATE settings SET
		 string_value = $2
		,int_value =    $3
		,json_value =   $4
		,time_value =   $5
		WHERE id =      $1`,
		model.ID,
		model.StringValue,
		model.IntValue,
		model.JSONValue,
		model.TimeValue,
	)
	return err
}
