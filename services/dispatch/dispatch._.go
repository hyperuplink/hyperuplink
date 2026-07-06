package dispatch

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
)

type Dispatch struct {
	cfg config.Redis
	ac  *asynq.Client
}

func New(cfg config.Redis) (disp *Dispatch, err error) {
	disp = new(Dispatch)

	disp.cfg = cfg
	disp.ac = nil

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
	if disp.ac != nil {
		disp.ac.Close()
	}
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
