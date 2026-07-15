package topics

import (
	"math"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type PollOption struct {
	Index    int    `json:"index"`
	Text     string `json:"text"`
	Votes    int    `json:"votes"`
	Percent  int    `json:"percent"`
	Selected bool   `json:"selected"`
}

type Poll struct {
	Options   []PollOption `json:"options"`
	Total     int          `json:"total"`
	Selection int          `json:"selection"`
	HasVoted  bool         `json:"has_voted"`
	Ended     bool         `json:"ended"`
	CanVote   bool         `json:"can_vote"`
}

type PollViewInput struct {
	Topic    *topic.Topic
	ViewerID uuid.NullUUID
	CanWrite bool
}

type PollVoteInput struct {
	Topic     *topic.Topic
	AuthorID  uuid.UUID
	Selection int
}

func PollView(
	rt *runtime.Runtime,
	in *PollViewInput,
) (poll *Poll, err error) {
	if !in.Topic.IsPoll() {
		return nil, nil
	}

	var tally map[int]int
	tally, err = rt.Repositories.PostEvent.TallyForTopicUUID(
		postevent.PollVote,
		in.Topic.ID,
		common.QueryOptions{},
	)
	if err != nil {
		return nil, err
	}

	poll = new(Poll)
	poll.Selection = -1
	poll.Ended = in.Topic.PollEnded()

	if in.ViewerID.Valid {
		var selection int
		var voted bool

		selection, voted, err = rt.Repositories.PostEvent.SelectionForTopicUUID(
			postevent.PollVote,
			in.Topic.ID,
			in.ViewerID.UUID,
			common.QueryOptions{Limit: 1},
		)
		if err != nil {
			return nil, err
		}

		poll.HasVoted = voted
		if voted {
			poll.Selection = selection
		}
	}

	for _, votes := range tally {
		poll.Total += votes
	}

	poll.Options = make([]PollOption, 0, len(in.Topic.PollOptions))
	for index, text := range in.Topic.PollOptions {
		poll.Options = append(poll.Options, PollOption{
			Index:    index,
			Text:     text,
			Votes:    tally[index],
			Percent:  pollPercent(tally[index], poll.Total),
			Selected: poll.HasVoted && poll.Selection == index,
		})
	}

	poll.CanVote = in.ViewerID.Valid &&
		in.CanWrite &&
		!poll.Ended &&
		!poll.HasVoted

	return poll, nil
}

func PollVote(
	rt *runtime.Runtime,
	in *PollVoteInput,
) (err error) {
	if !in.Topic.IsPoll() {
		return errs.ErrPollKindInvalid
	}

	if in.Topic.PollEnded() {
		return errs.ErrPollEnded
	}

	if !in.Topic.HasPollOption(in.Selection) {
		return errs.ErrPollSelectionInvalid
	}

	event := new(postevent.PostEvent)
	event.Type = postevent.PollVote
	event.AuthorID = in.AuthorID
	event.Target = postevent.Topic
	event.TopicID = uuid.NullUUID{UUID: in.Topic.ID, Valid: true}
	event.Selection = in.Selection

	_, err = rt.Repositories.PostEvent.Create(event)

	return err
}

func pollPercent(votes int, total int) int {
	if total <= 0 {
		return 0
	}

	return int(math.Round((float64(votes) / float64(total)) * 100))
}
