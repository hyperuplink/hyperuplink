package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	HashMem    = 64 * 1024
	HashIter   = 3
	HashSaltln = 16
	HashKeyln  = 32
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     Role      `json:"role"`

	Password            string           `json:"password"`
	PasswordResetToken  string           `json:"password_reset_token"`
	PasswordResetSentAt pgtype.Timestamp `json:"password_reset_sent_at"`

	Email                   string           `json:"email"`
	EmailUnconfirmed        string           `json:"email_unconfirmed"`
	EmailConfirmationToken  string           `json:"email_confirmation_token"`
	EmailConfirmationSentAt pgtype.Timestamp `json:"email_confirmation_sent_at"`
	EmailConfirmedAt        pgtype.Timestamp `json:"email_confirmed_at"`

	OTPEnabled  bool   `json:"otp_enabled"`
	OTPSecret   string `json:"otp_secret"`
	OTPTimestep int    `json:"otp_timestep"`

	SignInLastAt         pgtype.Timestamp `json:"sign_in_last_at"`
	SignInFailedAttempts int              `json:"sign_in_failed_attempts"`
	SignInLockedAt       pgtype.Timestamp `json:"sign_in_locked_at"`
	SignInUnlockToken    string           `json:"sign_in_unlock_token"`

	ProfilePicture string `json:"profile_picture"`
	Signature      string `json:"signature"`

	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	ConfirmedAt pgtype.Timestamp `json:"confirmed_at"`
	BannedAt    pgtype.Timestamp `json:"banned_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}
