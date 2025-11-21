package user

import (
	"time"

	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`

	Password            string    `json:"password"`
	PasswordResetToken  string    `json:"password_reset_token"`
	PasswordResetSentAt time.Time `json:"password_reset_sent_at"`

	Email                   string    `json:"email"`
	EmailUnconfirmed        string    `json:"email_unconfirmed"`
	EmailConfirmationToken  string    `json:"email_confirmation_token"`
	EmailConfirmationSentAt time.Time `json:"email_confirmation_sent_at"`
	EmailConfirmedAt        time.Time `json:"email_confirmed_at"`

	OTPEnabled  bool   `json:"otp_enabled"`
	OTPSecret   string `json:"otp_secret"`
	OTPTimestep int    `json:"otp_timestep"`

	SignInLastAt         time.Time `json:"sign_in_last_at"`
	SignInFailedAttempts int       `json:"sign_in_failed_attempts"`
	SignInLockedAt       time.Time `json:"sign_in_locked_at"`
	SignInUnlockToken    string    `json:"sign_in_unlock_token"`

	ProfilePicture string `json:"profile_picture"`
	Signature      string `json:"signature"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BannedAt  time.Time `json:"banned_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
