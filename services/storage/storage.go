package storage

import (
	"strings"

	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/services/config"
	"github.com/mrusme/hyperuplink/services/storage/s3"
)

type IStorage interface {
	Startup() error
	Shutdown() error
	StoreFile(src string, dest string) error
	GetFileDownloadURL(dest string) (string, error)
}

type Storage struct {
	cfg       config.Storages
	providers map[string]IStorage
}

func New(cfg config.Storages) (str *Storage, err error) {
	str = new(Storage)

	str.cfg = cfg
	str.providers = make(map[string]IStorage)

	for _, cfg := range str.cfg {
		switch strings.ToLower(cfg.Type) {
		case "s3":
			if str.providers[cfg.ID], err = s3.New(cfg); err != nil {
				return nil, err
			}
		default:
			return nil, errs.ErrStorageTypeInvalid
		}
	}

	return str, nil
}

func (str *Storage) Startup() (err error) {
	for _, provider := range str.providers {
		if err = provider.Startup(); err != nil {
			return err
		}
	}
	return nil
}

func (str *Storage) Shutdown() (err error) {
	for _, provider := range str.providers {
		if err = provider.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}

func (str *Storage) StoreFile(
	providerID string,
	src string,
	dest string,
) (err error) {
	if _, ok := str.providers[providerID]; !ok {
		return errs.ErrStorageIDNotFound
	}

	return str.providers[providerID].StoreFile(src, dest)
}

func (str *Storage) GetFileDownloadURL(
	providerID string,
	dest string,
) (dlurl string, err error) {
	if _, ok := str.providers[providerID]; !ok {
		return dlurl, errs.ErrStorageIDNotFound
	}

	return str.providers[providerID].GetFileDownloadURL(dest)
}
