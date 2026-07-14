package dispatch

import (
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/cron/cronrun"
)

const CRON_RETENTION time.Duration = 1 * time.Hour

func (disp *Dispatch) CronRun(payload *cronrun.CronRun) (err error) {
	j := asyncjob.New(
		payload.FunctionID,
		asyncjob.Cron,
		asyncjob.Run,
	)
	if err = j.SetPayload(payload); err != nil {
		return err
	}

	err = disp.enqueue(j, TaskCron,
		asynq.Queue(QueueCron),
		asynq.TaskID(payload.TaskID()),
		asynq.Retention(CRON_RETENTION),
	)
	if errors.Is(err, asynq.ErrTaskIDConflict) {
		return nil
	}

	return err
}
