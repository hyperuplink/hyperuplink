package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/worker/targets"
)

type Worker struct {
	rt       *runtime.Runtime
	ts       *targets.Targets
	redisCfg config.Redis
	as       *asynq.Server
	asMux    *asynq.ServeMux
}

func New(
	rt *runtime.Runtime,
) (*Worker, error) {
	wrk := new(Worker)

	wrk.rt = rt

	return wrk, nil
}

func (wrk *Worker) Startup() error {
	return nil
}

func (wrk *Worker) Run() (err error) {
	wrk.ts, err = targets.New(wrk.rt)
	wrk.rt.NilOrDie(err)

	err = wrk.ts.LoadAll()
	wrk.rt.NilOrDie(err)

	err = wrk.ts.RunAll()
	wrk.rt.NilOrDie(err)

	if wrk.redisCfg, err = wrk.rt.Config.Redis(); err != nil {
		return err
	}
	addrsl := len(wrk.redisCfg.Addrs)

	asyncConfig := asynq.Config{
		Logger: wrk.rt.ALogger,
		// Concurrency: wrk.redisCfg.Poolsize,
	}

	if addrsl == 1 {
		if wrk.rt.Config.RedisMasterName() == "" {
			wrk.as = asynq.NewServer(
				asynq.RedisClientOpt{
					Addr:     wrk.redisCfg.Addrs[0],
					DB:       wrk.redisCfg.Database,
					Username: wrk.redisCfg.Username,
					Password: wrk.redisCfg.Password,
					PoolSize: wrk.redisCfg.Poolsize,
				},
				asyncConfig,
			)
		} else {
			wrk.as = asynq.NewServer(
				asynq.RedisFailoverClientOpt{
					MasterName:    wrk.redisCfg.MasterName,
					SentinelAddrs: wrk.redisCfg.Addrs,
					DB:            wrk.redisCfg.Database,
					Username:      wrk.redisCfg.Username,
					Password:      wrk.redisCfg.Password,
					PoolSize:      wrk.redisCfg.Poolsize,
				},
				asyncConfig,
			)
		}
	} else {
		wrk.as = asynq.NewServer(
			asynq.RedisClusterClientOpt{
				Addrs:    wrk.redisCfg.Addrs,
				Username: wrk.redisCfg.Username,
				Password: wrk.redisCfg.Password,
			},
			asyncConfig,
		)
	}

	wrk.asMux = asynq.NewServeMux()
	wrk.asMux.HandleFunc("job", asynqHandler(wrk))

	err = wrk.as.Run(wrk.asMux)
	wrk.rt.NilOrDie(err)

	return err
}

func (wrk *Worker) Shutdown() error {
	if wrk.as != nil {
		wrk.as.Shutdown()
	}

	if wrk.ts == nil {
		return nil
	}

	ok, shutdownErrs := wrk.ts.ShutdownAll()
	if !ok {
		for id, err := range shutdownErrs {
			wrk.rt.Error("failed to shutdown target", "worker",
				"target", id, "error", err)
		}
	}

	return nil
}

func asynqHandler(wrk *Worker) func(context.Context, *asynq.Task) error {
	return wrk.HandleJob
}

func (wrk *Worker) HandleJob(ctx context.Context, t *asynq.Task) error {
	var job asyncjob.AsyncJob
	if err := json.Unmarshal(t.Payload(), &job); err != nil {
		return err
	}

	wrk.rt.Debug("status", "working", "payload", t.Payload())

	if err := wrk.ts.Execute(
		job.TargetID,
		job,
	); err != nil {
		wrk.rt.Error("error", err)
		return err
	}

	return nil
}
