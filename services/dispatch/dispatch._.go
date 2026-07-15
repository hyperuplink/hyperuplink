package dispatch

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/common"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const (
	TaskJob   string = "job"
	TaskCron  string = "cron"
	QueueCron string = "cron"
)

type Dispatch struct {
	cfg     config.Redis
	targets config.Targets
	repo    *repoSetting.Repository
	ac      *asynq.Client
}

func New(
	cfg config.Redis,
	targets config.Targets,
	repo *repoSetting.Repository,
) (disp *Dispatch, err error) {
	disp = new(Dispatch)

	disp.cfg = cfg
	disp.targets = targets
	disp.repo = repo
	disp.ac = nil

	return disp, nil
}

type routing struct {
	emailTargetID string
	xmppTargetID  string
}

func (r routing) targetIDFor(rcpt *common.Recipient) string {
	if rcpt.IsJID {
		return r.xmppTargetID
	}

	return r.emailTargetID
}

func (disp *Dispatch) debugTargetIDFor(channel string) string {
	for _, target := range disp.targets {
		if target.IsDebug() && target.Serves(channel) {
			return target.ID
		}
	}

	return ""
}

func (disp *Dispatch) routing() (r routing, err error) {
	settingCommsEmail, err := repoSetting.GetByID[setting.CommsEmail](
		disp.repo, "comms_email")
	if err != nil {
		return r, err
	}

	settingCommsXMPP, err := repoSetting.GetByID[setting.CommsXMPP](
		disp.repo, "comms_xmpp")
	if err != nil {
		return r, err
	}

	r.emailTargetID = settingCommsEmail.JSONValue.TargetID
	if r.emailTargetID == "" {
		r.emailTargetID = disp.debugTargetIDFor(config.TargetTypeEmail)
	}

	r.xmppTargetID = settingCommsXMPP.JSONValue.TargetID
	if r.xmppTargetID == "" {
		r.xmppTargetID = disp.debugTargetIDFor(config.TargetTypeXMPP)
	}

	return r, nil
}

func (disp *Dispatch) system() (sys *setting.System, err error) {
	settingSystem, err := repoSetting.GetByID[setting.System](
		disp.repo, "system")
	if err != nil {
		return nil, err
	}

	return &settingSystem.JSONValue, nil
}

func (disp *Dispatch) Startup() (err error) {
	addrsl := len(disp.cfg.Addrs)

	if addrsl == 1 {
		if disp.cfg.MasterName == "" {
			disp.ac = asynq.NewClient(asynq.RedisClientOpt{
				Addr:     disp.cfg.Addrs[0],
				DB:       disp.cfg.Database,
				Username: disp.cfg.Username,
				Password: disp.cfg.Password,
				PoolSize: disp.cfg.Poolsize,
			})
		} else {
			disp.ac = asynq.NewClient(asynq.RedisFailoverClientOpt{
				MasterName:    disp.cfg.MasterName,
				SentinelAddrs: disp.cfg.Addrs,
				DB:            disp.cfg.Database,
				Username:      disp.cfg.Username,
				Password:      disp.cfg.Password,
				PoolSize:      disp.cfg.Poolsize,
			})
		}
	} else {
		disp.ac = asynq.NewClient(asynq.RedisClusterClientOpt{
			Addrs:    disp.cfg.Addrs,
			Username: disp.cfg.Username,
			Password: disp.cfg.Password,
		})
	}

	return nil
}

func (disp *Dispatch) Shutdown() (err error) {
	if disp.ac != nil {
		disp.ac.Close()
	}
	return nil
}

func (disp *Dispatch) Job(j *asyncjob.AsyncJob) (err error) {
	return disp.enqueue(j, TaskJob)
}

func (disp *Dispatch) enqueue(
	j *asyncjob.AsyncJob,
	taskType string,
	opts ...asynq.Option,
) (err error) {
	if _, err = j.SetID(); err != nil {
		return err
	}

	jj, err := json.Marshal(j)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, jj, append([]asynq.Option{
		asynq.MaxRetry(5),
		asynq.Timeout(30 * time.Minute),
	}, opts...)...)

	_, err = disp.ac.Enqueue(task)

	return err
}
