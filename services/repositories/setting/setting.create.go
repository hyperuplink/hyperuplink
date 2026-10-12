package setting

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func Create[T any](
	repo *Repository,
	model *setting.Setting[T],
) (rowID string, err error) {
	return database.Value[string](context.Background(), repo.db, `INSERT INTO settings (
  id
  ,string_value
  ,int_value
	,json_value
  ,time_value
	) VALUES (
		$1
		,$2
		,$3
		,$4
		,$5
	) RETURNING id`,
		model.ID,
		model.StringValue,
		model.IntValue,
		model.JSONValue,
		model.TimeValue,
	)
}
