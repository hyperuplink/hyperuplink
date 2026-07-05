package site

import (
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/services/repositories/common"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
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
		s.r.GetRuntime().Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return nil
	}
	attachments := settingAttachments.JSONValue

	views = make([]AttachmentView, 0, len(ids))
	for _, id := range ids {
		model, merr := s.r.GetRuntime().Repositories.Attachment.GetByUUID(
			id,
			common.QueryOptions{},
		)
		if merr != nil {
			continue
		}

		dlurl, abs, uerr := s.r.GetRuntime().Storage.GetFileDownloadURL(
			attachments.StorageProviderID,
			path.Join(attachments.StoragePath, id.String()),
		)
		if uerr != nil {
			continue
		}
		if !abs {
			dlurl = s.HrefTo(strings.TrimPrefix(dlurl, "/"))
		}

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
