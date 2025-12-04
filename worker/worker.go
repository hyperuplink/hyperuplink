package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/runtime"
	"github.com/mrusme/hyperuplink/worker/targets"
)

type Worker struct {
	rt       *runtime.Runtime
	ts       *targets.Targets
	redis    *asynq.Server
	redisMux *asynq.ServeMux
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

	addrs := wrk.rt.Config.RedisAddresses()
	addrsl := len(addrs)

	asyncConfig := asynq.Config{
		Logger:      wrk.rt.ALogger,
		Concurrency: wrk.rt.Config.RedisPoolsize(),
	}

	if addrsl == 1 {
		if wrk.rt.Config.RedisMasterName() == "" {
			wrk.redis = asynq.NewServer(
				asynq.RedisClientOpt{
					Addr:     addrs[0],
					Username: wrk.rt.Config.RedisUsername(),
					Password: wrk.rt.Config.RedisPassword(),
				},
				asyncConfig,
			)
		} else {
			wrk.redis = asynq.NewServer(
				asynq.RedisFailoverClientOpt{
					MasterName:    wrk.rt.Config.RedisMasterName(),
					SentinelAddrs: addrs,
					Username:      wrk.rt.Config.RedisUsername(),
					Password:      wrk.rt.Config.RedisPassword(),
				},
				asyncConfig,
			)
		}
	} else {
		wrk.redis = asynq.NewServer(
			asynq.RedisClusterClientOpt{
				Addrs:    addrs,
				Username: wrk.rt.Config.RedisUsername(),
				Password: wrk.rt.Config.RedisPassword(),
			},
			asyncConfig,
		)
	}

	wrk.redisMux = asynq.NewServeMux()
	wrk.redisMux.HandleFunc("message", asynqHandler(wrk))

	err = wrk.redis.Run(wrk.redisMux)
	wrk.rt.NilOrDie(err)

	return err
}

func (wrk *Worker) Shutdown() error {
	wrk.redis.Shutdown()
	return nil
}

func asynqHandler(wrk *Worker) func(context.Context, *asynq.Task) error {
	return wrk.HandleJob
}

func (wrk *Worker) HandleJob(ctx context.Context, t *asynq.Task) error {
	var m asyncjob.AsyncJob
	if err := json.Unmarshal(t.Payload(), &m); err != nil {
		return err
	}

	wrk.rt.Debug("status", "working", "payload", t.Payload())

	return nil
}
