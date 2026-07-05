package helpers

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/attachment"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/runtime"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

const AttachmentsFormField = "attachments"

func ProcessAttachments(
	rt *runtime.Runtime,
	c fiber.Ctx,
	authorID uuid.UUID,
) (ids []uuid.UUID, err error) {
	var settingAttachments *setting.Setting[setting.Attachments]
	settingAttachments, err = settingRepo.GetByID[setting.Attachments](
		rt.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return nil, err
	}
	attachments := settingAttachments.JSONValue

	if !attachments.EnableAttachments {
		return nil, nil
	}

	form, ferr := c.MultipartForm()
	if ferr != nil {
		return nil, nil
	}

	files := form.File[AttachmentsFormField]
	if len(files) == 0 {
		return nil, nil
	}

	ids = make([]uuid.UUID, 0, len(files))
	for _, fh := range files {
		var id uuid.UUID
		if id, err = processAttachment(rt, &attachments, authorID, fh); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func processAttachment(
	rt *runtime.Runtime,
	attachments *setting.Attachments,
	authorID uuid.UUID,
	fh *multipart.FileHeader,
) (id uuid.UUID, err error) {
	if fh.Size > attachments.GetMaxSize() {
		return uuid.Nil, errors.New("attachment_too_large")
	}

	src, err := fh.Open()
	if err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "hyperuplink-attachment-*")
	if err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err = io.Copy(tmp, src); err != nil {
		tmp.Close()
		return uuid.Nil, errors.New("attachment_upload_failed")
	}
	if err = tmp.Close(); err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}

	mtype, err := mimetype.DetectFile(tmpPath)
	if err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}

	allowed := false
	for _, format := range attachments.UploadFormats {
		if mtype.Is(format) {
			allowed = true
			break
		}
	}
	if !allowed {
		return uuid.Nil, errors.New("attachment_format_not_allowed")
	}

	output, herr := attachments.RunOnUploadHook(tmpPath)
	if herr != nil {
		return uuid.Nil, errors.New("attachment_hook_failed")
	}

	var checksum string
	if checksum, err = checksumFile(tmpPath); err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
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
			return uuid.Nil, errors.New("attachment_duplicate")
		}
		return uuid.Nil, errors.New("attachment_upload_failed")
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}
	defer f.Close()

	if err = rt.Storage.StoreFile(
		attachments.StorageProviderID,
		f,
		path.Join(attachments.StoragePath, id.String()),
	); err != nil {
		return uuid.Nil, errors.New("attachment_upload_failed")
	}

	return id, nil
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
