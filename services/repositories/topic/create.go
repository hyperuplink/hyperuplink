package topic

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/topic"
)

func (repo *Repository) Create(model *topic.Topic) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO topics (
		 name
		,slug
		,forum_id
		,author_id
		,kind
		,anonymous
		,text
		,poll_options
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
		,NOW()
		,NOW()
		,$9
	) RETURNING id`,
		model.Name,
		model.Slug,
		model.ForumID,
		model.AuthorID,
		model.Kind,
		model.Anonymous,
		model.Text,
		model.PollOptions,
		model.EndedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
