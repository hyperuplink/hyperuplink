package site

import (
	"strings"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type AttachmentView struct {
	ID       string
	Filename string
	URL      string
	AuthorID string
	Inline   bool
}

func (s *Site) Attachments(ids []uuid.UUID) (views []AttachmentView) {
	if len(ids) == 0 {
		return nil
	}

	var settingAttachments *setting.Setting[setting.Attachments]
	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		gh.Repositories(s.r.GetRuntime()).Setting,
		"attachments",
	)
	if err != nil {
		return nil
	}
	attachments := settingAttachments.JSONValue

	views = make([]AttachmentView, 0, len(ids))
	for _, id := range ids {
		model, merr := gh.Repositories(s.r.GetRuntime()).Attachment.GetByUUID(
			id,
			common.QueryOptions{},
		)
		if merr != nil {
			continue
		}

		dlurl := s.HrefTo(route.For("Attachment").Fill(map[string]string{
			"attachment": id.String(),
		}).AsURL())

		views = append(views, AttachmentView{
			ID:       id.String(),
			Filename: model.Filename,
			URL:      dlurl,
			AuthorID: model.AuthorID.String(),
			Inline: attachments.InlineImageDisplay &&
				strings.HasPrefix(model.MimeType, "image/"),
		})
	}

	return views
}
