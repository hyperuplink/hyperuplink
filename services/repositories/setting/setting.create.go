package setting

import (
	"github.com/mrusme/hyperuplink/models/setting"
)

func Create[T any](
	repo *Repository,
	model *setting.Setting[T],
) (rowID string, err error) {
	err = repo.db.QueryRow(`INSERT INTO settings (
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
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
