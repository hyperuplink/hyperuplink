package config

type Storage struct {
	ID    string       `koanf:"ID"`
	Type  string       `koanf:"Type"`
	S3    StorageS3    `koanf:"S3,omitempty"`
	Local StorageLocal `koanf:"Local,omitempty"`
}

type StorageS3 struct {
	Endpoint          string
	Region            string
	AccessKey         string
	SecretKey         string
	PublicURL         string
	PublicDownload    bool
	PresignedDownload bool
}

type StorageLocal struct {
	Path      string
	PublicURI string
}

type Storages []Storage

func (cfg *Config) Storages() (storages Storages, err error) {
	err = cfg.k.Unmarshal("Storage", &storages)
	return storages, err
}
