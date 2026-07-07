package permission

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	None      byte = 0b000
	ReadOnly  byte = 0b100
	ReadWrite byte = 0b110
	Moderate  byte = 0b111
)

const (
	LevelNone              = "none"
	LevelRead              = "read"
	LevelReadWrite         = "read_write"
	LevelReadWriteModerate = "read_write_moderate"
)

type Permission struct {
	ID         uuid.UUID   `json:"id"`
	GroupID    pgtype.Text `json:"group_id"`
	CategoryID pgtype.UUID `json:"category_id"`
	Bits       pgtype.Bits `json:"bits"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}

func (m *Permission) SetLevel(level byte) {
	m.Bits.Bytes = []byte{level << 5}
	m.Bits.Len = 3
	m.Bits.Valid = true
}

func (m *Permission) Level() byte {
	if m.Bits.Valid == false || len(m.Bits.Bytes) == 0 {
		return None
	}
	return m.Bits.Bytes[0] >> 5
}

func (m *Permission) IsNone() bool {
	return m.Level() == None
}

func (m *Permission) IsRead() bool {
	return m.Level() == ReadOnly
}

func (m *Permission) IsReadWrite() bool {
	return m.Level() == ReadWrite
}

func (m *Permission) IsModerate() bool {
	return m.Level() == Moderate
}

func (m *Permission) LevelString() string {
	return LevelToString(m.Level())
}

func LevelFromString(level string) (byte, bool) {
	switch level {
	case LevelNone:
		return None, true
	case LevelRead:
		return ReadOnly, true
	case LevelReadWrite:
		return ReadWrite, true
	case LevelReadWriteModerate:
		return Moderate, true
	}
	return None, false
}

func LevelToString(level byte) string {
	switch level {
	case ReadOnly:
		return LevelRead
	case ReadWrite:
		return LevelReadWrite
	case Moderate:
		return LevelReadWriteModerate
	}
	return LevelNone
}

func LevelToBitString(level byte) string {
	switch level {
	case ReadOnly:
		return "100"
	case ReadWrite:
		return "110"
	case Moderate:
		return "111"
	}
	return "000"
}
