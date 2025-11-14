package repositories

import (
	"github.com/mrusme/hyperuplink/services/database"
	"github.com/mrusme/hyperuplink/services/repositories/user"
)

type Repositories struct {
	db   database.IDatabase
	User *user.Repository
	// Application *application.Repository
	// Target      *target.Repository
}

func New(
	db database.IDatabase,
) (*Repositories, error) {
	var repos *Repositories = new(Repositories)
	var err error

	repos.db = db

	var userRepo *user.Repository
	if userRepo, err = user.New(repos.db); err != nil {
		return nil, err
	}
	//
	// var appRepo *application.Repository
	// if appRepo, err = application.New(cfg, db); err != nil {
	// 	return nil, err
	// }
	//
	// var targetRepo *target.Repository
	// if targetRepo, err = target.New(cfg, db); err != nil {
	// 	return nil, err
	// }

	repos.User = userRepo
	// repos.Application = appRepo
	// repos.Target = targetRepo

	return repos, nil
}

func (repos *Repositories) Startup() error {
	var err error

	if err = repos.User.Startup(); err != nil {
		return err
	}

	return nil
}

func (repos *Repositories) Shutdown() error {
	return nil
}
