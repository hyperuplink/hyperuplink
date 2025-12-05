package dispatch

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/notification/replynotification"
	"github.com/mrusme/hyperuplink/services/config"
)

type Dispatch struct {
	cfg config.Redis
	ac  *asynq.Client
}

func New(cfg config.Redis) (disp *Dispatch, err error) {
	disp = new(Dispatch)

	disp.cfg = cfg

	return disp, nil
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
	disp.ac.Close()
	return nil
}

func (disp *Dispatch) Job(j *asyncjob.AsyncJob) (err error) {
	if _, err = j.SetID(); err != nil {
		return err
	}

	jj, err := json.Marshal(j)
	if err != nil {
		return err
	}

	task := asynq.NewTask("job", jj,
		asynq.MaxRetry(5), asynq.Timeout(30*time.Minute))

	_, err = disp.ac.Enqueue(task)

	return err
}

func (disp *Dispatch) ReplyNotifications(
	targetID string,
	payload []replynotification.ReplyNotification,
) (err error) {
	j := asyncjob.New(
		targetID,
		asyncjob.Notification,
		asyncjob.Reply,
	)
	if err = j.SetPayload(payload); err != nil {
		return err
	}

	if err = disp.Job(j); err != nil {
		return err
	}

	return nil
}

func (disp *Dispatch) ReplyNotification(
	targetID string,
	payload replynotification.ReplyNotification,
) (err error) {
	return disp.ReplyNotifications(
		targetID,
		[]replynotification.ReplyNotification{payload},
	)
}
