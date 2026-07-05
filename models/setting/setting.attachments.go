package setting

import (
	"os/exec"
	"strings"
)

const DEFAULT_ATTACHMENT_MAX_SIZE int64 = 10485760

const AttachmentHookPlaceholder = "%ATTACHMENT%"

var AttachmentUploadFormatOptions = []string{
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/webp",
	"application/json",
	"application/pdf",
	"application/zip",
	"application/gzip",
	"audio/mpeg",
	"audio/ogg",
	"audio/wav",
	"video/mp4",
	"video/webm",
	"text/plain",
	"text/csv",
	"text/markdown",
	"text/html",
}

type Attachments struct {
	EnableAttachments bool     `json:"enable_attachments"`
	UploadFormats     []string `json:"upload_formats"` // Options: image/gif, image/jpeg, image/png, image/webp, application/json, application/pdf, application/zip, audio/mpeg, audio/vorbis, video/mp4, text/plain, text/csv, text/html
	MaxSize           int64    `json:"max_size"`       // in byte
	StorageProviderID string   `json:"storage_provider_id"`
	StoragePath       string   `json:"storage_path"`
	OnUploadHook      string   `json:"on_upload_hook"`
}

func (a *Attachments) GetMaxSize() (size int64) {
	if a.MaxSize <= 0 {
		a.MaxSize = DEFAULT_ATTACHMENT_MAX_SIZE
	}

	return a.MaxSize
}

func (a *Attachments) RunOnUploadHook(attachmentPath string) (err error) {
	hook := strings.TrimSpace(a.OnUploadHook)
	if hook == "" {
		return nil
	}

	fields := strings.Fields(hook)
	for i, field := range fields {
		fields[i] = strings.ReplaceAll(field, AttachmentHookPlaceholder, attachmentPath)
	}

	cmd := exec.Command(fields[0], fields[1:]...)
	return cmd.Run()
}
