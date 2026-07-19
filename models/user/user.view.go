package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Public struct {
	ID             uuid.UUID        `json:"id"`
	Username       string           `json:"username"`
	Role           Role             `json:"role"`
	MemberOf       []string         `json:"member_of"`
	ProfilePicture string           `json:"profile_picture"`
	SignatureText  string           `json:"signature_text"`
	SignatureHTML  string           `json:"signature_html"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
}

type Detail struct {
	Public

	Email            string           `json:"email"`
	EmailUnconfirmed string           `json:"email_unconfirmed"`
	EmailIsJID       bool             `json:"email_is_jid"`
	Language         string           `json:"language"`
	Timezone         string           `json:"timezone"`
	OTPEnabled       bool             `json:"otp_enabled"`
	SignInLastAt     pgtype.Timestamp `json:"sign_in_last_at"`
	ConfirmedAt      pgtype.Timestamp `json:"confirmed_at"`
	BannedAt         pgtype.Timestamp `json:"banned_at"`
	DeletedAt        pgtype.Timestamp `json:"deleted_at"`
}

func (m *User) AsPublic() Public {
	return Public{
		ID:             m.ID,
		Username:       m.Username,
		Role:           m.Role,
		MemberOf:       m.MemberOf,
		ProfilePicture: m.ProfilePicture,
		SignatureText:  m.SignatureText,
		SignatureHTML:  m.SignatureHTML,
		CreatedAt:      m.CreatedAt,
	}
}

func (m *User) AsDetail() Detail {
	return Detail{
		Public:           m.AsPublic(),
		Email:            m.Email,
		EmailUnconfirmed: m.EmailUnconfirmed,
		EmailIsJID:       m.EmailIsJID,
		Language:         m.Language,
		Timezone:         m.Timezone,
		OTPEnabled:       m.OTPEnabled,
		SignInLastAt:     m.SignInLastAt,
		ConfirmedAt:      m.ConfirmedAt,
		BannedAt:         m.BannedAt,
		DeletedAt:        m.DeletedAt,
	}
}
