package cron

import (
	"fmt"
	"sync"

	robfig "github.com/robfig/cron/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
)

type Fn func() error

type Entry struct {
	ID       string
	Spec     string
	Schedule robfig.Schedule
	Fn       Fn
}

type Cron struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

func New() (crn *Cron, err error) {
	crn = new(Cron)
	crn.entries = make(map[string]*Entry)

	return crn, nil
}

func (crn *Cron) Startup() (err error) {
	return nil
}

func (crn *Cron) Shutdown() (err error) {
	return nil
}

func (crn *Cron) Register(id string, spec string, fn Fn) (err error) {
	if id == "" || fn == nil {
		return errs.ErrCronFunctionInvalid
	}

	var schedule robfig.Schedule
	if schedule, err = robfig.ParseStandard(spec); err != nil {
		return fmt.Errorf("%w: %s: %w", errs.ErrCronFunctionInvalid, id, err)
	}

	crn.mu.Lock()
	defer crn.mu.Unlock()

	if _, ok := crn.entries[id]; ok {
		return fmt.Errorf("%w: %s", errs.ErrCronFunctionIDExists, id)
	}

	crn.entries[id] = &Entry{
		ID:       id,
		Spec:     spec,
		Schedule: schedule,
		Fn:       fn,
	}

	return nil
}

func (crn *Cron) Get(id string) (entry *Entry, ok bool) {
	crn.mu.RLock()
	defer crn.mu.RUnlock()

	entry, ok = crn.entries[id]

	return entry, ok
}

func (crn *Cron) Entries() (entries []*Entry) {
	crn.mu.RLock()
	defer crn.mu.RUnlock()

	entries = make([]*Entry, 0, len(crn.entries))
	for _, entry := range crn.entries {
		entries = append(entries, entry)
	}

	return entries
}
