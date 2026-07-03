package setting

type Profiles struct {
	EnablePicture            bool   `json:"enable_picture"`
	PictureFormat            string `json:"picture_format"` // Options: webp, png, jpg
	PictureStorageProviderID string `json:"picture_storage_provider_id"`
	PictureStoragePath       string `json:"picture_storage_path"`
}
