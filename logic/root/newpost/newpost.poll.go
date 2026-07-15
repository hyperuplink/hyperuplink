package newpost

import (
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const (
	PollOptionsMin    int = 2
	PollOptionsMax    int = 8
	PollOptionMaxRune int = 78
)

var pollEndsAtLayouts = []string{
	"2006-01-02T15:04",
	"2006-01-02T15:04:05",
}

type PollInput struct {
	Options  []string
	EndsAt   string
	Location *time.Location
}

func PollAllowed(rt *runtime.Runtime) (allowed bool, err error) {
	var settingTopics *setting.Setting[setting.Topics]

	if settingTopics, err = repoSetting.GetByID[setting.Topics](
		rt.Repositories.Setting,
		"topics",
	); err != nil {
		return false, err
	}

	return settingTopics.JSONValue.AllowKindPoll, nil
}

func ApplyPoll(
	rt *runtime.Runtime,
	top *topic.Topic,
	in *PollInput,
) (err error) {
	var allowed bool
	if allowed, err = PollAllowed(rt); err != nil {
		return err
	}
	if !allowed {
		return errs.ErrPollNotAllowed
	}

	options := TrimPollOptions(in.Options)

	if len(options) < PollOptionsMin {
		return errs.ErrPollOptionsTooFew
	}
	if len(options) > PollOptionsMax {
		return errs.ErrPollOptionsTooMany
	}
	for _, option := range options {
		if len([]rune(option)) > PollOptionMaxRune {
			return errs.ErrPollOptionTooLong
		}
	}

	var endedAt pgtype.Timestamp
	if endedAt, err = ParsePollEndsAt(in.EndsAt, in.Location); err != nil {
		return err
	}

	top.Kind = topic.Poll
	top.PollOptions = options
	top.EndedAt = endedAt

	return nil
}

func TrimPollOptions(options []string) (trimmed []string) {
	for _, option := range options {
		option = strings.TrimSpace(option)
		if option == "" {
			continue
		}
		trimmed = append(trimmed, option)
	}

	return trimmed
}

func ParsePollEndsAt(
	value string,
	loc *time.Location,
) (endedAt pgtype.Timestamp, err error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return pgtype.Timestamp{}, nil
	}

	if loc == nil {
		loc = time.UTC
	}

	var parsed time.Time
	var ok bool
	for _, layout := range pollEndsAtLayouts {
		if parsed, err = time.ParseInLocation(layout, value, loc); err == nil {
			ok = true
			break
		}
	}
	if !ok {
		return pgtype.Timestamp{}, errs.ErrPollEndsAtInvalid
	}

	if !parsed.After(time.Now()) {
		return pgtype.Timestamp{}, errs.ErrPollEndsAtPast
	}

	return pgtype.Timestamp{Time: parsed.UTC(), Valid: true}, nil
}
