package session

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type Session struct {
	c fiber.Ctx
}

func New(c fiber.Ctx) *Session {
	s := new(Session)
	s.c = c

	return s
}

func (s *Session) GetID() string {
	sess := session.FromContext(s.c)
	return sess.ID()
}

func (s *Session) Set(authProvider string, userId string) error {
	sess := session.FromContext(s.c)

	if err := sess.Regenerate(); err != nil {
		return err
	}

	sess.Set("auth_provider", authProvider)
	sess.Set("user_id", userId)

	return nil
}

func (s *Session) GetProvider() (string, bool) {
	sess := session.FromContext(s.c)

	authProvider := sess.Get("auth_provider")
	if authProvider == nil || authProvider.(string) == "" {
		return "", false
	}

	return authProvider.(string), true
}

func (s *Session) GetUserID() (string, bool) {
	sess := session.FromContext(s.c)

	userID := sess.Get("user_id")

	if userID == nil || userID.(string) == "" {
		return "", false
	}

	return userID.(string), true
}

func (s *Session) Reset() error {
	sess := session.FromContext(s.c)

	return sess.Reset()
}
