package activity

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	activityRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/activity"
)

const (
	DEFAULT_FLUSH_INTERVAL time.Duration = 5 * time.Second
	DEFAULT_MAX_BUFFERED   int           = 8192
)

type Activity struct {
	log  *slog.Logger
	repo *activityRepo.Repository

	interval time.Duration
	maxBuf   int

	mu        sync.Mutex
	buffered  map[activity.Key]*activity.Record
	running   bool
	stop      chan struct{}
	stopped   chan struct{}
	flushSoon chan struct{}
}

func New(
	log *slog.Logger,
	repo *activityRepo.Repository,
) (act *Activity, err error) {
	act = new(Activity)

	act.log = log
	act.repo = repo
	act.interval = DEFAULT_FLUSH_INTERVAL
	act.maxBuf = DEFAULT_MAX_BUFFERED
	act.buffered = make(map[activity.Key]*activity.Record)

	return act, nil
}

func (act *Activity) Startup() (err error) {
	act.mu.Lock()
	defer act.mu.Unlock()

	if act.running {
		return nil
	}

	act.stop = make(chan struct{})
	act.stopped = make(chan struct{})
	act.flushSoon = make(chan struct{}, 1)
	act.running = true

	go act.run()

	return nil
}

func (act *Activity) Shutdown() (err error) {
	act.mu.Lock()
	if !act.running {
		act.mu.Unlock()
		return nil
	}
	act.running = false
	stopped := act.stopped
	close(act.stop)
	act.mu.Unlock()

	<-stopped

	return nil
}

func (act *Activity) Record(rec activity.Record) (err error) {
	if rec.Count < 1 {
		rec.Count = 1
	}

	if rec.Kind.Policy() == activity.Immediate {
		return act.repo.Append([]activity.Record{rec})
	}

	act.mu.Lock()
	if cur, ok := act.buffered[rec.Key]; ok {
		cur.Count += rec.Count
	} else {
		buf := rec
		act.buffered[rec.Key] = &buf
	}
	full := len(act.buffered) >= act.maxBuf
	act.mu.Unlock()

	if full {
		select {
		case act.flushSoon <- struct{}{}:
		default:
		}
	}

	return nil
}

func (act *Activity) drain() (recs []activity.Record) {
	act.mu.Lock()
	if len(act.buffered) == 0 {
		act.mu.Unlock()
		return nil
	}
	buffered := act.buffered
	act.buffered = make(map[activity.Key]*activity.Record)
	act.mu.Unlock()

	recs = make([]activity.Record, 0, len(buffered))
	for _, rec := range buffered {
		recs = append(recs, *rec)
	}

	return recs
}

func (act *Activity) Flush() (err error) {
	recs := act.drain()
	if len(recs) == 0 {
		return nil
	}

	var coalescing []activity.Record
	var appending []activity.Record

	for _, rec := range recs {
		if rec.Kind.Coalesces() {
			coalescing = append(coalescing, rec)
		} else {
			appending = append(appending, rec)
		}
	}

	return errors.Join(
		act.repo.Upsert(coalescing),
		act.repo.Append(appending),
	)
}

func (act *Activity) flushLogged() {
	if err := act.Flush(); err != nil {
		act.log.Error("Activity.Flush", "error", err)
	}
}

func (act *Activity) run() {
	defer close(act.stopped)

	ticker := time.NewTicker(act.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			act.flushLogged()
		case <-act.flushSoon:
			act.flushLogged()
		case <-act.stop:
			act.flushLogged()
			return
		}
	}
}
