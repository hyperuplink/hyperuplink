package setting

const DEFAULT_PICTURE_MAX_SIZE int64 = 1048576

var PictureUploadFormatOptions = []string{
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/webp",
}

type Profiles struct {
	EnablePicture            bool     `json:"enable_picture"`
	PictureUploadFormats     []string `json:"upload_formats"`   // Options: image/gif, image/jpeg, image/png, image/webp
	PictureFormat            string   `json:"picture_format"`   // Options: webp, png, jpg
	PictureMaxSize           int64    `json:"picture_max_size"` // in byte
	PictureStorageProviderID string   `json:"picture_storage_provider_id"`
	PictureStoragePath       string   `json:"picture_storage_path"`
}

func (p *Profiles) GetPictureMaxSize() (size int64) {
	if p.PictureMaxSize <= 0 {
		p.PictureMaxSize = DEFAULT_PICTURE_MAX_SIZE
	}

	return p.PictureMaxSize
}
