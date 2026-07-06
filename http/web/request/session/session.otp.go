package session

import "github.com/gofiber/fiber/v3/middleware/session"

const pendingOTPURLKey = "otp_pending_url"

func (s *Session) SetPendingOTPURL(url string) {
	sess := session.FromContext(s.c)
	sess.Set(pendingOTPURLKey, url)
}

func (s *Session) GetPendingOTPURL() (string, bool) {
	sess := session.FromContext(s.c)

	val := sess.Get(pendingOTPURLKey)
	if val == nil {
		return "", false
	}

	url, ok := val.(string)
	if !ok || url == "" {
		return "", false
	}

	return url, true
}

func (s *Session) ClearPendingOTPURL() {
	sess := session.FromContext(s.c)
	sess.Delete(pendingOTPURLKey)
}
