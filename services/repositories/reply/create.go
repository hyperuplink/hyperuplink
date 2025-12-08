package reply

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/reply"
)

func (repo *Repository) Create(model *reply.Reply) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO replies (
		 id
		,topic_id
		,reply_id
		,author_id
		,kind
		,text
		,poll_vote
		,rsvp_response`+
		// ,created_at
		// ,updated_at
		// ,moderated_at
		// ,spammed_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,$8`+
		// ,$9
		// ,$10
		// ,$11
		// ,$12
		// ,$13
		`) RETURNING id`,
		model.ID,
		model.TopicID,
		model.ReplyID,
		model.AuthorID,
		model.Kind,
		model.Text,
		model.PollVote,
		model.RSVPResponse,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.ModeratedAt,
		// model.SpammedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
