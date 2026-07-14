package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/cron/cronrun"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/services/dispatch"
)

const CHECK_INTERVAL time.Duration = 5 * time.Second

type Cron struct {
	rt       *runtime.Runtime
	redisCfg config.Redis
	as       *asynq.Server
	asMux    *asynq.ServeMux

	next    map[string]time.Time
	stop    chan struct{}
	stopped chan struct{}
}

func New(rt *runtime.Runtime) (crn *Cron, err error) {
	crn = new(Cron)

	crn.rt = rt
	crn.next = make(map[string]time.Time)
	crn.stop = make(chan struct{})
	crn.stopped = make(chan struct{})

	return crn, nil
}

func (crn *Cron) Startup() (err error) {
	return nil
}

func (crn *Cron) Run() (err error) {
	if crn.redisCfg, err = crn.rt.Config.Redis(); err != nil {
		return err
	}

	asyncConfig := asynq.Config{
		Logger: crn.rt.ALogger,
		Queues: map[string]int{dispatch.QueueCron: 1},
	}

	addrsl := len(crn.redisCfg.Addrs)

	if addrsl == 1 {
		if crn.redisCfg.MasterName == "" {
			crn.as = asynq.NewServer(
				asynq.RedisClientOpt{
					Addr:     crn.redisCfg.Addrs[0],
					DB:       crn.redisCfg.Database,
					Username: crn.redisCfg.Username,
					Password: crn.redisCfg.Password,
					PoolSize: crn.redisCfg.Poolsize,
				},
				asyncConfig,
			)
		} else {
			crn.as = asynq.NewServer(
				asynq.RedisFailoverClientOpt{
					MasterName:    crn.redisCfg.MasterName,
					SentinelAddrs: crn.redisCfg.Addrs,
					DB:            crn.redisCfg.Database,
					Username:      crn.redisCfg.Username,
					Password:      crn.redisCfg.Password,
					PoolSize:      crn.redisCfg.Poolsize,
				},
				asyncConfig,
			)
		}
	} else {
		crn.as = asynq.NewServer(
			asynq.RedisClusterClientOpt{
				Addrs:    crn.redisCfg.Addrs,
				Username: crn.redisCfg.Username,
				Password: crn.redisCfg.Password,
			},
			asyncConfig,
		)
	}

	crn.asMux = asynq.NewServeMux()
	crn.asMux.HandleFunc(dispatch.TaskCron, crn.HandleJob)

	if err = crn.as.Start(crn.asMux); err != nil {
		crn.rt.Error("error", err)
		return err
	}

	crn.rt.Info("status", "ok", "entries", len(crn.rt.Cron.Entries()))

	crn.schedule()

	return nil
}

func (crn *Cron) Shutdown() (err error) {
	select {
	case <-crn.stop:
	default:
		close(crn.stop)
	}

	<-crn.stopped

	if crn.as != nil {
		crn.as.Shutdown()
	}

	return nil
}

func (crn *Cron) schedule() {
	defer close(crn.stopped)

	ticker := time.NewTicker(CHECK_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			crn.tick(now)
		case <-crn.stop:
			return
		}
	}
}

func (crn *Cron) tick(now time.Time) {
	for _, entry := range crn.rt.Cron.Entries() {
		next, ok := crn.next[entry.ID]
		if !ok {
			crn.next[entry.ID] = entry.Schedule.Next(now)
			continue
		}

		if now.Before(next) {
			continue
		}

		crn.dispatch(entry.ID, next)
		crn.next[entry.ID] = entry.Schedule.Next(now)
	}
}

func (crn *Cron) dispatch(id string, scheduledAt time.Time) {
	crn.rt.Debug("dispatch", id, "scheduled_at", scheduledAt)

	if err := crn.rt.Dispatch.CronRun(
		cronrun.New(id, scheduledAt),
	); err != nil {
		crn.rt.Error("error", err)
	}
}

func (crn *Cron) HandleJob(ctx context.Context, t *asynq.Task) (err error) {
	var job asyncjob.AsyncJob
	if err = json.Unmarshal(t.Payload(), &job); err != nil {
		return err
	}

	if job.Type != asyncjob.Cron || job.SubType != asyncjob.Run {
		return fmt.Errorf("%w: %s/%s",
			errs.ErrJobPayloadInvalid, job.Type, job.SubType)
	}

	var payload cronrun.CronRun
	if err = json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}

	entry, ok := crn.rt.Cron.Get(payload.FunctionID)
	if !ok {
		return fmt.Errorf("%w: %s",
			errs.ErrCronFunctionIDNotFound, payload.FunctionID)
	}

	crn.rt.Info("status", "running", "function_id", entry.ID)

	if err = entry.Fn(); err != nil {
		crn.rt.Error("error", err, "function_id", entry.ID)
		return err
	}

	crn.rt.Info("status", "ok", "function_id", entry.ID)

	return nil
}
