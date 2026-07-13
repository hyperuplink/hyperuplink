package runtime

import (
	"fmt"
	"time"
)

const ResetTimeLayout = "15:04"

func ValidateResetTime(token string, now time.Time) error {
	t, err := time.Parse(ResetTimeLayout, token)
	if err != nil {
		return fmt.Errorf(
			"invalid --reset confirmation %q: expected the current time as HH:MM in 24h format, e.g. %s",
			token, now.Format(ResetTimeLayout),
		)
	}

	if t.Hour() != now.Hour() || t.Minute() != now.Minute() {
		return fmt.Errorf(
			"--reset confirmation %s does not match the current system time %s; refusing to reset",
			token, now.Format(ResetTimeLayout),
		)
	}

	return nil
}

func (rt *Runtime) Reset() (err error) {
	rt.Debug("status", "exec")
	rt.Warn("reset", "clearing database")

	if err = rt.Database.Reset(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Info("status", "ok", "reset", "database cleared")
	return nil
}
