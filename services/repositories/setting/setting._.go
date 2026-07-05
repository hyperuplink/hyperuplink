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
			Name:        "Hyperuplink",
			BaseURL:     "http://localhost:3000",
			Theme:       setting.DEFAULT_THEME,
			Colorscheme: setting.DEFAULT_COLORSCHEME,
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

	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}
