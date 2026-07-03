package setting

import (
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/services/database"
)

type Repository struct {
	db *database.Database
}

func New(db *database.Database) (*Repository, error) {
	repo := new(Repository)
	repo.db = db

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

	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}
