package session

import (
	"strconv"

	"github.com/gofiber/fiber/v3/middleware/session"
)

const (
	pending2FAProviderKey = "pending_2fa_provider"
	pending2FAUserIDKey   = "pending_2fa_user_id"
	pending2FAAttemptsKey = "pending_2fa_attempts"
)

func (s *Session) SetPending2FA(authProvider, userID string) {
	sess := session.FromContext(s.c)
	sess.Set(pending2FAProviderKey, authProvider)
	sess.Set(pending2FAUserIDKey, userID)
	sess.Delete(pending2FAAttemptsKey)
}

func (s *Session) GetPending2FA() (userID string, authProvider string, ok bool) {
	sess := session.FromContext(s.c)

	uid, uok := sess.Get(pending2FAUserIDKey).(string)
	prov, pok := sess.Get(pending2FAProviderKey).(string)
	if !uok || !pok || uid == "" || prov == "" {
		return "", "", false
	}

	return uid, prov, true
}

func (s *Session) ClearPending2FA() {
	sess := session.FromContext(s.c)
	sess.Delete(pending2FAProviderKey)
	sess.Delete(pending2FAUserIDKey)
	sess.Delete(pending2FAAttemptsKey)
}

func (s *Session) IncrementPending2FAAttempts() int {
	sess := session.FromContext(s.c)

	var n int
	if v, ok := sess.Get(pending2FAAttemptsKey).(string); ok {
		n, _ = strconv.Atoi(v)
	}
	n++
	sess.Set(pending2FAAttemptsKey, strconv.Itoa(n))

	return n
}
