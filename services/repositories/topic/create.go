package topic

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/topic"
)

func (repo *Repository) Create(model *topic.Topic) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO topics (
		 id
		,name
		,slug
		,forum_id
		,author_id
		,kind
		,anonymous
		,text
		,poll_options`+
		// ,created_at
		// ,updated_at
		`,ended_at`+
		// ,moderated_at
		// ,spammed_at
		// ,locked_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,$8
		,$9
		,$10`+
		// ,$11
		// ,$12
		// ,$13
		// ,$14
		// ,$15
		// ,$16
		`) RETURNING id`,
		model.ID,
		model.Name,
		model.Slug,
		model.ForumID,
		model.AuthorID,
		model.Kind,
		model.Anonymous,
		model.Text,
		model.PollOptions,
		// model.CreatedAt,
		// model.UpdatedAt,
		model.EndedAt,
		// model.ModeratedAt,
		// model.SpammedAt,
		// model.LockedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
