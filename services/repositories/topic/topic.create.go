package topic

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/topic"
)

func (repo *Repository) Create(model *topic.Topic) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO topics (
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
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
