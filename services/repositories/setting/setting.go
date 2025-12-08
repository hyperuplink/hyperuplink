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
		settingSystem.JSONValue = setting.System{
			Name:    "Hyperuplink",
			BaseURL: "http://localhost:3000",
		}
		if _, err = Create(repo, settingSystem); err != nil {
			return err
		}
	}

	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}
