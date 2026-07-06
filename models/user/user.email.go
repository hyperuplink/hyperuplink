package user

import (
	"time"

	"xn--gckvb8fzb.com/hyperuplink/errs"
)

func (m *User) SetEmailConfirmationSentAt(t time.Time) {
	m.EmailConfirmationSentAt.Time = t
	m.EmailConfirmationSentAt.Valid = !t.IsZero()
}

func (m *User) SetEmailConfirmedAt(t time.Time) {
	m.EmailConfirmedAt.Time = t
	m.EmailConfirmedAt.Valid = !t.IsZero()
}

func (m *User) ClearEmailConfirmedAt() {
	m.EmailConfirmedAt.Time = time.Time{}
	m.EmailConfirmedAt.Valid = false
}

func (m *User) ResetEmailConfirmationToken() error {
	const charset string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
	const length uint32 = 8

	b, err := m.generateRandomBytes(length)
	if err != nil {
		return err
	}

	for i := uint32(0); i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	m.EmailConfirmationToken = string(b)

	return nil
}

func (m *User) SetEmailForConfirmation(email string) (err error) {
	m.EmailUnconfirmed = email
	if err = m.ResetEmailConfirmationToken(); err != nil {
		return err
	}
	m.ClearEmailConfirmedAt()
	m.SetEmailConfirmationSentAt(time.Now())

	return nil
}

func (m *User) ConfirmEmail(token string) error {
	now := time.Now()

	if m.EmailConfirmationToken != token {
		return errs.ErrEmailConfirmationTokenWrong
	}

	m.Email = m.EmailUnconfirmed
	m.EmailUnconfirmed = ""
	m.EmailConfirmationToken = ""
	m.SetEmailConfirmedAt(now)

	if m.ConfirmedAt.Valid == false || m.ConfirmedAt.Time.IsZero() {
		m.SetConfirmedAt(now)
	}

	return nil
}

func (m *User) SetEmailIsJID() {
	m.EmailIsJID = true
}
