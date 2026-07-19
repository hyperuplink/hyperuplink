package attachments

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

type UpdateInput struct {
	EnableAttachments  bool     `json:"enable_attachments" form:"enable_attachments"`
	UploadFormats      []string `json:"upload_formats" form:"upload_formats" validate:"required_if=EnableAttachments true"`
	MaxSize            int64    `json:"max_size" form:"max_size" validate:"required,min=1"`
	StorageProviderID  string   `json:"storage_provider_id" form:"storage_provider_id" validate:"required_if=EnableAttachments true,max=64"`
	StoragePath        string   `json:"storage_path" form:"storage_path" validate:"omitempty,max=255"`
	OnUploadHook       string   `json:"on_upload_hook" form:"on_upload_hook" validate:"omitempty,max=1024"`
	InlineImageDisplay bool     `json:"inline_image_display" form:"inline_image_display"`
}

type View struct {
	Attachments         *setting.Attachments `json:"attachments"`
	StorageIDs          []string             `json:"storage_ids"`
	StorageIsPublic     bool                 `json:"storage_is_public"`
	UploadFormatOptions []string             `json:"upload_format_options"`
}
