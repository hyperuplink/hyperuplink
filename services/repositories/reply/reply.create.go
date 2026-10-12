package reply

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
)

func (repo *Repository) Create(ctx context.Context, model *reply.Reply) (rowID uuid.UUID, err error) {
	return database.Value[uuid.UUID](ctx, repo.db, `INSERT INTO replies (
		short_id
		,topic_id
		,reply_id
		,author_id
		,text
		,html
		,attachment_ids
		,created_at
		,updated_at
	) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,NOW()
		,NOW()
		) RETURNING id`,
		model.ShortID,
		model.TopicID,
		model.ReplyID,
		model.AuthorID,
		model.Text,
		model.HTML,
		model.AttachmentIDs,
	)
}
