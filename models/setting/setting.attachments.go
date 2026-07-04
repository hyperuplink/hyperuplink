package setting

type Attachments struct {
	EnableAttachments bool     `json:"enable_attachments"`
	UploadFormats     []string `json:"upload_formats"` // Options: image/gif, image/jpeg, image/png, image/webp, application/json, application/pdf, application/zip, audio/mpeg, audio/vorbis, video/mp4, text/plain, text/csv, text/html
	MaxSize           int64    `json:"max_size"`       // in byte
	StorageProviderID string   `json:"storage_provider_id"`
	StoragePath       string   `json:"storage_path"`
	OnUploadHook      string   `json:"on_upload_hook"`
}
