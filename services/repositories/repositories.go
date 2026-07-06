package repositories

import (
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/services/database"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/attachment"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/category"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/forum"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/permission"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/postevent"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/reply"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/search"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/topic"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/unit"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/user"
)

type Repositories struct {
	db         *database.Database
	Setting    *setting.Repository
	Unit       *unit.Repository
	User       *user.Repository
	Category   *category.Repository
	Forum      *forum.Repository
	Topic      *topic.Repository
	Reply      *reply.Repository
	Attachment *attachment.Repository
	Permission *permission.Repository
	PostEvent  *postevent.Repository
	Search     *search.Repository
}

func New(
	db *database.Database,
	cfg *config.Config,
) (repos *Repositories, err error) {
	repos = new(Repositories)
	repos.db = db

	var settingRepo *setting.Repository
	if settingRepo, err = setting.New(repos.db, cfg); err != nil {
		return nil, err
	}
	repos.Setting = settingRepo

	var unitRepo *unit.Repository
	if unitRepo, err = unit.New(repos.db); err != nil {
		return nil, err
	}
	repos.Unit = unitRepo

	var userRepo *user.Repository
	if userRepo, err = user.New(repos.db); err != nil {
		return nil, err
	}
	repos.User = userRepo

	var categoryRepo *category.Repository
	if categoryRepo, err = category.New(repos.db); err != nil {
		return nil, err
	}
	repos.Category = categoryRepo

	var forumRepo *forum.Repository
	if forumRepo, err = forum.New(repos.db); err != nil {
		return nil, err
	}
	repos.Forum = forumRepo

	var topicRepo *topic.Repository
	if topicRepo, err = topic.New(repos.db); err != nil {
		return nil, err
	}
	repos.Topic = topicRepo

	var replyRepo *reply.Repository
	if replyRepo, err = reply.New(repos.db); err != nil {
		return nil, err
	}
	repos.Reply = replyRepo

	var attachmentRepo *attachment.Repository
	if attachmentRepo, err = attachment.New(repos.db); err != nil {
		return nil, err
	}
	repos.Attachment = attachmentRepo

	var permissionRepo *permission.Repository
	if permissionRepo, err = permission.New(repos.db); err != nil {
		return nil, err
	}
	repos.Permission = permissionRepo

	var posteventRepo *postevent.Repository
	if posteventRepo, err = postevent.New(repos.db); err != nil {
		return nil, err
	}
	repos.PostEvent = posteventRepo

	var searchRepo *search.Repository
	if searchRepo, err = search.New(repos.db); err != nil {
		return nil, err
	}
	repos.Search = searchRepo

	return repos, nil
}

func (repos *Repositories) Startup() (err error) {
	if err = repos.Setting.Startup(); err != nil {
		return err
	}

	if err = repos.Unit.Startup(); err != nil {
		return err
	}

	if err = repos.User.Startup(); err != nil {
		return err
	}

	if err = repos.Category.Startup(); err != nil {
		return err
	}

	if err = repos.Forum.Startup(); err != nil {
		return err
	}

	if err = repos.Topic.Startup(); err != nil {
		return err
	}

	if err = repos.Reply.Startup(); err != nil {
		return err
	}

	if err = repos.Attachment.Startup(); err != nil {
		return err
	}

	if err = repos.Permission.Startup(); err != nil {
		return err
	}

	if err = repos.PostEvent.Startup(); err != nil {
		return err
	}

	if err = repos.Search.Startup(); err != nil {
		return err
	}

	return nil
}

func (repos *Repositories) Shutdown() (err error) {
	if err = repos.Search.Shutdown(); err != nil {
		return err
	}

	if err = repos.PostEvent.Shutdown(); err != nil {
		return err
	}

	if err = repos.Permission.Shutdown(); err != nil {
		return err
	}

	if err = repos.Attachment.Shutdown(); err != nil {
		return err
	}

	if err = repos.Reply.Shutdown(); err != nil {
		return err
	}

	if err = repos.Topic.Shutdown(); err != nil {
		return err
	}

	if err = repos.Forum.Shutdown(); err != nil {
		return err
	}

	if err = repos.Category.Shutdown(); err != nil {
		return err
	}

	if err = repos.User.Shutdown(); err != nil {
		return err
	}

	if err = repos.Unit.Shutdown(); err != nil {
		return err
	}

	if err = repos.Setting.Shutdown(); err != nil {
		return err
	}

	return nil
}
