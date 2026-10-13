package memo

import (
	"sync"
	"time"
)

const DEFAULT_TTL time.Duration = time.Minute

type Memo struct {
	mu         sync.Mutex
	ttl        time.Duration
	generation uint64
	entries    map[string]entry
}

type entry struct {
	value   any
	expires time.Time
}

func New(ttl time.Duration) *Memo {
	if ttl <= 0 {
		ttl = DEFAULT_TTL
	}

	return &Memo{
		ttl:     ttl,
		entries: map[string]entry{},
	}
}

func (m *Memo) Reset() {
	m.mu.Lock()
	m.generation++
	clear(m.entries)
	m.mu.Unlock()
}

func Get[T any](m *Memo, key string, load func() (T, error)) (value T, err error) {
	now := time.Now()

	m.mu.Lock()
	generation := m.generation
	if e, ok := m.entries[key]; ok && now.Before(e.expires) {
		if value, ok = e.value.(T); ok {
			m.mu.Unlock()
			return value, nil
		}
	}
	m.mu.Unlock()

	if value, err = load(); err != nil {
		return value, err
	}

	m.mu.Lock()
	if generation == m.generation {
		for k, e := range m.entries {
			if !now.Before(e.expires) {
				delete(m.entries, k)
			}
		}
		m.entries[key] = entry{value: value, expires: now.Add(m.ttl)}
	}
	m.mu.Unlock()

	return value, nil
}
