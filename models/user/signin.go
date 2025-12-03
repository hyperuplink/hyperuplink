package user

import "time"

func (m *User) SetSignInLastAt(t time.Time) {
	m.SignInLastAt.Time = t
	m.SignInLastAt.Valid = !t.IsZero()
}

func (m *User) SetSignInLockedAt(t time.Time) {
	m.SignInLockedAt.Time = t
	m.SignInLockedAt.Valid = !t.IsZero()
}
