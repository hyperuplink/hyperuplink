package user

import (
	"time"

	"xn--gckvb8fzb.com/glides/helpers/password"
)

type PasswordParams = password.Params

func (m *User) generateRandomBytes(n uint32) ([]byte, error) {
	return password.RandomBytes(n)
}

func (m *User) SetRandomPassword() (err error) {
	m.Password, err = password.HashRandom()
	return err
}

func (m *User) SetPassword(pw string) (err error) {
	m.Password, err = password.Hash(pw)
	return err
}

func (m *User) CheckPassword(pw string) (
	match bool,
	params *PasswordParams,
	err error,
) {
	return password.Check(m.Password, pw)
}

func (m *User) SetPasswordResetSentAt(t time.Time) {
	m.PasswordResetSentAt.Time = t
	m.PasswordResetSentAt.Valid = !t.IsZero()
}
