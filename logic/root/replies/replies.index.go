package replies

import (
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type IndexInput struct {
	Topic   *vtopic.VTopic
	Page    int
	PerPage int
}

type Page struct {
	Replies []vreply.VReply
	Total   int64
}

func Index(
	rt *runtime.Runtime,
	in *IndexInput,
) (page *Page, err error) {
	page = new(Page)

	if page.Replies, page.Total, err = gh.Repositories(rt).Reply.VAllForTopicUUIDOrdered(
		in.Topic.ID,
		in.PerPage,
		in.Page,
	); err != nil {
		return nil, err
	}

	return page, nil
}
