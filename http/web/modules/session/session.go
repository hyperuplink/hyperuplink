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

func (s *Session) GetProvider() string {
	sess := session.FromContext(s.c)

	return sess.Get("auth_provider").(string)
}

func (s *Session) Reset() error {
	sess := session.FromContext(s.c)

	return sess.Reset()
}
