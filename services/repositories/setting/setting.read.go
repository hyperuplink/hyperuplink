package setting

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/memo"
)

type row struct {
	ID          string
	StringValue string
	IntValue    int
	JSONValue   []byte
	TimeValue   pgtype.Timestamp
}

func GetByID[T any](
	ctx context.Context,
	repo *Repository,
	id string,
) (model *setting.Setting[T], err error) {
	stored, err := memo.Get(repo.memo, "setting:"+id, func() (row, error) {
		loaded, lerr := database.Row[row](ctx, repo.db, `SELECT * FROM settings
			WHERE id = $1
			LIMIT 1`,
			id,
		)
		if lerr != nil {
			return row{}, lerr
		}

		return *loaded, nil
	})
	if err != nil {
		return nil, err
	}

	model = &setting.Setting[T]{
		ID:          stored.ID,
		StringValue: stored.StringValue,
		IntValue:    stored.IntValue,
		TimeValue:   stored.TimeValue,
	}
	if len(stored.JSONValue) > 0 {
		if err = json.Unmarshal(stored.JSONValue, &model.JSONValue); err != nil {
			return nil, err
		}
	}

	return model, nil
}
