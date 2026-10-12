package topic

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

func (repo *Repository) Create(ctx context.Context, model *topic.Topic) (rowID uuid.UUID, err error) {
	return database.Value[uuid.UUID](ctx, repo.db, `INSERT INTO topics (
		short_id
		,name
		,slug
		,forum_id
		,author_id
		,kind
		,anonymous
		,text
		,html
		,poll_options
		,attachment_ids
		,created_at
		,updated_at
		,ended_at
		) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,$8
		,$9
		,$10
		,$11
		,NOW()
		,NOW()
		,$12
	) RETURNING id`,
		model.ShortID,
		model.Name,
		model.Slug,
		model.ForumID,
		model.AuthorID,
		model.Kind,
		model.Anonymous,
		model.Text,
		model.HTML,
		model.PollOptions,
		model.AttachmentIDs,
		// CreatedAt
		// UpdatedAt
		model.EndedAt,
	)
}
