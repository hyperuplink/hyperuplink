package user

import "time"

func (m *User) SetEmailConfirmationSentAt(t time.Time) {
	m.EmailConfirmationSentAt.Time = t
	m.EmailConfirmationSentAt.Valid = !t.IsZero()
}

func (m *User) SetEmailConfirmedAt(t time.Time) {
	m.EmailConfirmedAt.Time = t
	m.EmailConfirmedAt.Valid = !t.IsZero()
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
