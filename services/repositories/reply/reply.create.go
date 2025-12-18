package reply

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/reply"
)

func (repo *Repository) Create(model *reply.Reply) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO replies (
		topic_id
		,reply_id
		,author_id
		,kind
		,text
		,html
		,poll_vote
		,rsvp_response
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
		,$8
		,NOW()
		,NOW()
		) RETURNING id`,
		model.TopicID,
		model.ReplyID,
		model.AuthorID,
		model.Kind,
		model.Text,
		model.HTML,
		model.PollVote,
		model.RSVPResponse,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
