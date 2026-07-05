package setting

import (
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/services/config"
	"github.com/mrusme/hyperuplink/services/database"
)

type Repository struct {
	db  *database.Database
	cfg *config.Config
}

func New(db *database.Database, cfg *config.Config) (*Repository, error) {
	repo := new(Repository)
	repo.db = db
	repo.cfg = cfg

	return repo, nil
}

func (repo *Repository) Startup() (err error) {
	var settingSystem *setting.Setting[setting.System]

	if settingSystem, err = GetByID[setting.System](repo, "system"); err != nil {
		settingSystem = new(setting.Setting[setting.System])
		settingSystem.ID = "system"
		settingSystem.JSONValue = setting.System{
			Name:    "Hyperuplink",
			BaseURL: "http://localhost:3000",
		}
		if _, err = Create(repo, settingSystem); err != nil {
			return err
		}
	}

	var settingTopics *setting.Setting[setting.Topics]

	if settingTopics, err = GetByID[setting.Topics](repo, "topics"); err != nil {
		settingTopics = new(setting.Setting[setting.Topics])
		settingTopics.ID = "topics"
		settingTopics.JSONValue = setting.Topics{
			AllowKindQuestion: true,
			AllowKindPoll:     true,
			AllowKindRSVP:     true,
		}
		if _, err = Create(repo, settingTopics); err != nil {
			return err
		}
	}

	var settingProfiles *setting.Setting[setting.Profiles]

	if settingProfiles, err = GetByID[setting.Profiles](repo, "profiles"); err != nil {
		settingProfiles = new(setting.Setting[setting.Profiles])
		settingProfiles.ID = "profiles"
		settingProfiles.JSONValue = setting.Profiles{
			EnablePicture:            false,
			PictureUploadFormats:     setting.PictureUploadFormatOptions,
			PictureFormat:            "webp",
			PictureMaxSize:           setting.DEFAULT_PICTURE_MAX_SIZE,
			PictureStorageProviderID: "",
			PictureStoragePath:       "profile-pictures",
		}
		if _, err = Create(repo, settingProfiles); err != nil {
			return err
		}
	} else if settingProfiles.JSONValue.PictureStorageProviderID != "" {
		var storages config.Storages
		if storages, err = repo.cfg.Storages(); err != nil {
			return err
		}

		exists := false
		for _, storage := range storages {
			if storage.ID == settingProfiles.JSONValue.PictureStorageProviderID {
				exists = true
				break
			}
		}

		if !exists {
			settingProfiles.JSONValue.EnablePicture = false
			settingProfiles.JSONValue.PictureStorageProviderID = ""
			if err = Update(repo, settingProfiles); err != nil {
				return err
			}
		}
	}

	var settingAttachments *setting.Setting[setting.Attachments]

	if settingAttachments, err = GetByID[setting.Attachments](repo, "attachments"); err != nil {
		settingAttachments = new(setting.Setting[setting.Attachments])
		settingAttachments.ID = "attachments"
		settingAttachments.JSONValue = setting.Attachments{
			EnableAttachments: false,
			UploadFormats:     setting.AttachmentUploadFormatOptions,
			MaxSize:           setting.DEFAULT_ATTACHMENT_MAX_SIZE,
			StorageProviderID: "",
			StoragePath:       "attachments",
			OnUploadHook:      "",
		}
		if _, err = Create(repo, settingAttachments); err != nil {
			return err
		}
	} else if settingAttachments.JSONValue.StorageProviderID != "" {
		var storages config.Storages
		if storages, err = repo.cfg.Storages(); err != nil {
			return err
		}

		exists := false
		for _, storage := range storages {
			if storage.ID == settingAttachments.JSONValue.StorageProviderID {
				exists = true
				break
			}
		}

		if !exists {
			settingAttachments.JSONValue.EnableAttachments = false
			settingAttachments.JSONValue.StorageProviderID = ""
			if err = Update(repo, settingAttachments); err != nil {
				return err
			}
		}
	}

	var settingTheme *setting.Setting[setting.Theme]

	if settingTheme, err = GetByID[setting.Theme](repo, "theme"); err != nil {
		settingTheme = new(setting.Setting[setting.Theme])
		settingTheme.ID = "theme"
		settingTheme.JSONValue = setting.Theme{
			Theme:                  setting.DEFAULT_THEME,
			Colorscheme:            setting.DEFAULT_COLORSCHEME,
			ThemeStorageProviderID: "",
			ThemeStoragePath:       "theme",
			CustomBanner:           "",
			CustomFavicon:          "",
		}
		if _, err = Create(repo, settingTheme); err != nil {
			return err
		}
	} else if settingTheme.JSONValue.ThemeStorageProviderID != "" {
		var storages config.Storages
		if storages, err = repo.cfg.Storages(); err != nil {
			return err
		}

		exists := false
		for _, storage := range storages {
			if storage.ID == settingTheme.JSONValue.ThemeStorageProviderID {
				exists = true
				break
			}
		}

		if !exists {
			settingTheme.JSONValue.ThemeStorageProviderID = ""
			settingTheme.JSONValue.CustomBanner = ""
			if err = Update(repo, settingTheme); err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}
