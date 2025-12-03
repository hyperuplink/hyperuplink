package user

import "time"

func (m *User) SetCreatedAt(t time.Time) {
	m.CreatedAt.Time = t
	m.CreatedAt.Valid = !t.IsZero()
}

func (m *User) SetUpdatedAt(t time.Time) {
	m.UpdatedAt.Time = t
	m.UpdatedAt.Valid = !t.IsZero()
}

func (m *User) SetConfirmedAt(t time.Time) {
	m.ConfirmedAt.Time = t
	m.ConfirmedAt.Valid = !t.IsZero()
}

func (m *User) SetBannedAt(t time.Time) {
	m.BannedAt.Time = t
	m.BannedAt.Valid = !t.IsZero()
}

func (m *User) SetDeletedAt(t time.Time) {
	m.DeletedAt.Time = t
	m.DeletedAt.Valid = !t.IsZero()
}
