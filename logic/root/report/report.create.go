package report

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Create(
	rt *runtime.Runtime,
	authorID uuid.UUID,
	in *CreateInput,
) (err error) {
	post, err := ResolvePost(rt, in.Target, in.ID)
	if err != nil {
		return err
	}

	event := new(postevent.PostEvent)
	event.Type = postevent.Report
	event.AuthorID = authorID
	event.Target = post.Target
	event.TopicID = post.TopicID
	event.ReplyID = post.ReplyID
	event.Selection = in.ReportType

	_, err = rt.Repositories.PostEvent.Create(event)

	return err
}
