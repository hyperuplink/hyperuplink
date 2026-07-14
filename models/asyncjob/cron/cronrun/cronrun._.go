package cronrun

import (
	"fmt"
	"time"
)

type CronRun struct {
	FunctionID  string `json:"function_id"`
	ScheduledAt int64  `json:"scheduled_at"`
}

func New(functionID string, scheduledAt time.Time) (entity *CronRun) {
	entity = new(CronRun)
	entity.FunctionID = functionID
	entity.ScheduledAt = scheduledAt.Unix()

	return entity
}

func (entity *CronRun) TaskID() string {
	return fmt.Sprintf("cron:%s:%d", entity.FunctionID, entity.ScheduledAt)
}
