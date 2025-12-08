package repositories

import (
	"github.com/mrusme/hyperuplink/services/database"
	"github.com/mrusme/hyperuplink/services/repositories/setting"
	"github.com/mrusme/hyperuplink/services/repositories/user"
)

type Repositories struct {
	db      *database.Database
	Setting *setting.Repository
	User    *user.Repository
}

func New(
	db *database.Database,
) (repos *Repositories, err error) {
	repos = new(Repositories)
	repos.db = db

	var settingRepo *setting.Repository
	if settingRepo, err = setting.New(repos.db); err != nil {
		return nil, err
	}
	repos.Setting = settingRepo

	var userRepo *user.Repository
	if userRepo, err = user.New(repos.db); err != nil {
		return nil, err
	}
	repos.User = userRepo

	return repos, nil
}

func (repos *Repositories) Startup() (err error) {
	if err = repos.Setting.Startup(); err != nil {
		return err
	}

	if err = repos.User.Startup(); err != nil {
		return err
	}

	return nil
}

func (repos *Repositories) Shutdown() (err error) {
	if err = repos.User.Shutdown(); err != nil {
		return err
	}

	if err = repos.Setting.Shutdown(); err != nil {
		return err
	}

	return nil
}
