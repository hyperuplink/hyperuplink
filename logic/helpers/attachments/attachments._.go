package attachments

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/attachment"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"

	"github.com/gabriel-vasile/mimetype"
)

const FormField = "attachments"

func Settings(rt *runtime.Runtime) (*setting.Attachments, error) {
	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		rt.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return nil, err
	}

	return &settingAttachments.JSONValue, nil
}

func StoreAll(
	rt *runtime.Runtime,
	authorID uuid.UUID,
	files []*multipart.FileHeader,
) (ids []uuid.UUID, err error) {
	attachments, err := Settings(rt)
	if err != nil {
		return nil, err
	}

	if !attachments.EnableAttachments {
		return nil, nil
	}

	if len(files) == 0 {
		return nil, nil
	}

	ids = make([]uuid.UUID, 0, len(files))
	for _, fh := range files {
		var id uuid.UUID
		if id, err = Store(rt, attachments, authorID, fh); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func Store(
	rt *runtime.Runtime,
	attachments *setting.Attachments,
	authorID uuid.UUID,
	fh *multipart.FileHeader,
) (id uuid.UUID, err error) {
	if fh.Size > attachments.GetMaxSize() {
		return uuid.Nil, errs.ErrAttachmentTooLarge
	}

	src, err := fh.Open()
	if err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "hyperuplink-attachment-*")
	if err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err = io.Copy(tmp, src); err != nil {
		tmp.Close()
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}
	if err = tmp.Close(); err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}

	mtype, err := mimetype.DetectFile(tmpPath)
	if err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}

	allowed := false
	for _, format := range attachments.UploadFormats {
		if mtype.Is(format) {
			allowed = true
			break
		}
	}
	if !allowed {
		return uuid.Nil, errs.ErrAttachmentFormatNotAllowed
	}

	output, herr := attachments.RunOnUploadHook(tmpPath)
	if herr != nil {
		return uuid.Nil, errs.ErrAttachmentHookFailed
	}

	var checksum string
	if checksum, err = checksumFile(tmpPath); err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}

	att := &attachment.Attachment{
		AuthorID:           authorID,
		Filename:           fh.Filename,
		MimeType:           mtype.String(),
		Checksum:           checksum,
		OnUploadHookOutput: output,
	}
	if id, err = rt.Repositories.Attachment.Create(att); err != nil {
		if strings.HasPrefix(err.Error(), "unique_violation") {
			return uuid.Nil, errs.ErrAttachmentDuplicate
		}
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}
	defer f.Close()

	if err = rt.Storage.StoreFile(
		attachments.StorageProviderID,
		f,
		path.Join(attachments.StoragePath, id.String()),
	); err != nil {
		return uuid.Nil, errs.ErrAttachmentUploadFailed
	}

	return id, nil
}

func ResolveOwned(
	rt *runtime.Runtime,
	authorID uuid.UUID,
	ids []string,
) (owned []uuid.UUID, err error) {
	owned = make([]uuid.UUID, 0, len(ids))
	for _, raw := range ids {
		var id uuid.UUID
		if id, err = uuid.Parse(raw); err != nil {
			return nil, errs.ErrAttachmentNotFound
		}

		att, aerr := rt.Repositories.Attachment.GetByUUID(id, common.QueryOptions{})
		if aerr != nil {
			return nil, errs.ErrAttachmentNotFound
		}
		if att.AuthorID != authorID {
			return nil, errs.ErrAttachmentNotFound
		}

		owned = append(owned, id)
	}

	return owned, nil
}

func checksumFile(p string) (sum string, err error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
