package forum

import "time"

func (m *Forum) SetCreatedAt(t time.Time) {
	m.CreatedAt.Time = t
	m.CreatedAt.Valid = !t.IsZero()
}

func (m *Forum) SetUpdatedAt(t time.Time) {
	m.UpdatedAt.Time = t
	m.UpdatedAt.Valid = !t.IsZero()
}

func (m *Forum) SetDeletedAt(t time.Time) {
	m.DeletedAt.Time = t
	m.DeletedAt.Valid = !t.IsZero()
}
