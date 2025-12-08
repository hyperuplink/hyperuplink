package permission

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/models/user"
)

const (
	ReadOnly  byte = 0b100
	ReadWrite      = 0b110
	Moderate       = 0b111
)

type Permission struct {
	ID      uuid.UUID   `json:"id"`
	Role    user.Role   `json:"role"`
	Unit    pgtype.Text `json:"unit"`
	ForumID uuid.UUID   `json:"forum_id"`
	Bits    pgtype.Bits `json:"bits"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}

func (m *Permission) SetReadOnly() {
	m.Bits.Bytes = []byte{ReadOnly}
	m.Bits.Valid = true
}

func (m *Permission) SetReadWrite() {
	m.Bits.Bytes = []byte{ReadWrite}
	m.Bits.Valid = true
}

func (m *Permission) SetModerate() {
	m.Bits.Bytes = []byte{Moderate}
	m.Bits.Valid = true
}

func (m *Permission) IsReadOnly() bool {
	if m.Bits.Valid == true &&
		len(m.Bits.Bytes) > 0 &&
		m.Bits.Bytes[0] == ReadOnly {
		return true
	}

	return false
}

func (m *Permission) IsReadWrite() bool {
	if m.Bits.Valid == true &&
		len(m.Bits.Bytes) > 0 &&
		m.Bits.Bytes[0] == ReadWrite {
		return true
	}

	return false
}

func (m *Permission) IsModerate() bool {
	if m.Bits.Valid == true &&
		len(m.Bits.Bytes) > 0 &&
		m.Bits.Bytes[0] == Moderate {
		return true
	}

	return false
}
