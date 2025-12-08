package setting

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Setting[T any] struct {
	ID          string           `json:"id"`
	StringValue string           `json:"string_value"`
	IntValue    int              `json:"int_value"`
	JSONValue   T                `json:"json_value"`
	TimeValue   pgtype.Timestamp `json:"time_value"`
}

