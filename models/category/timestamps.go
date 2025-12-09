package category

import "time"

func (m *Category) SetCreatedAt(t time.Time) {
	m.CreatedAt.Time = t
	m.CreatedAt.Valid = !t.IsZero()
}

func (m *Category) SetUpdatedAt(t time.Time) {
	m.UpdatedAt.Time = t
	m.UpdatedAt.Valid = !t.IsZero()
}

func (m *Category) SetDeletedAt(t time.Time) {
	m.DeletedAt.Time = t
	m.DeletedAt.Valid = !t.IsZero()
}
