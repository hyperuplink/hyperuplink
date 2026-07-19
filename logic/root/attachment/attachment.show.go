package attachment

import (
	"io"
	"path"
	"strings"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type File struct {
	Attachment  *attachment.Attachment
	Reader      io.Reader
	Disposition string
}

var inlineImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

func Open(
	rt *runtime.Runtime,
	id uuid.UUID,
	perms *permission.Resolution,
) (file *File, err error) {
	att, err := rt.Repositories.Attachment.GetByUUID(id, common.QueryOptions{})
	if err != nil {
		return nil, errs.ErrNoRows
	}

	categoryID, found, err := rt.Repositories.Attachment.GetCategoryForAttachment(id)
	if err != nil || !found {
		return nil, errs.ErrNoRows
	}

	if !perms.CanReadID(categoryID) {
		return nil, errs.ErrNoRows
	}

	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		rt.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return nil, errs.ErrNoRows
	}
	attachments := settingAttachments.JSONValue

	reader, err := rt.Storage.GetFile(
		attachments.StorageProviderID,
		path.Join(attachments.StoragePath, id.String()),
	)
	if err != nil {
		return nil, errs.ErrNoRows
	}

	// Serving user-uploaded content like image/svg+xml or text/html inline
	// would allow scripts to execute and access the session cookie. Therefor
	// only serve images inline.
	disposition := "attachment"
	if attachments.InlineImageDisplay && inlineImageTypes[att.MimeType] {
		disposition = "inline"
	}

	return &File{
		Attachment:  att,
		Reader:      reader,
		Disposition: disposition,
	}, nil
}

func SanitizeFilename(name string) string {
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == '"' || r == '\\' || r == '/' {
			return -1
		}
		return r
	}, name)
	if name == "" {
		return "attachment"
	}
	return name
}
