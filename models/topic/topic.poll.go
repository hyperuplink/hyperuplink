package topic

import (
	"time"
)

func (m Topic) IsPoll() bool {
	return m.Kind == Poll
}

func (m Topic) PollEnded() bool {
	return m.EndedAt.Valid && m.EndedAt.Time.Before(time.Now())
}

func (m Topic) HasPollOption(index int) bool {
	return index >= 0 && index < len(m.PollOptions)
}
