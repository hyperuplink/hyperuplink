package postevent

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
)

func (repo *Repository) Create(model *postevent.PostEvent) (rowID uuid.UUID, err error) {
	return database.Value[uuid.UUID](context.Background(), repo.db, `INSERT INTO postevents (
		 type
		,author_id
		,target
		,topic_id
		,reply_id
		,selection
		,created_at
		,updated_at
	) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,$6
		,NOW()
		,NOW()
	) RETURNING id`,
		model.Type,
		model.AuthorID,
		model.Target,
		model.TopicID,
		model.ReplyID,
		model.Selection,
	)
}
