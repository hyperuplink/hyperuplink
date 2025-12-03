package user

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) Create(model *user.User) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO users (
  username
  ,role
  ,password`+
		// ,password_reset_token
		// ,password_reset_sent_at
		`,email
  ,email_unconfirmed
  ,email_confirmation_token
  ,email_confirmation_sent_at`+
		// ,email_confirmed_at
		// ,otp_enabled
		// ,otp_secret
		// ,otp_timestep
		// ,sign_in_last_at
		// ,sign_in_failed_attempts
		// ,sign_in_locked_at
		// ,sign_in_unlock_token
		// ,profile_picture
		// ,signature
		// ,created_at
		// ,updated_at
		// ,confirmed_at
		// ,banned_at
		// ,deleted_at
		`) VALUES (
	$1
	,$2
	,$3
		,$4
		,$5
		,$6
	,$7`+
		// ,$8
		// ,$9
		// ,$10
		// ,$11
		// ,$12
		// ,$13
		// ,$14
		// ,$15
		// ,$16
		// ,$17
		// ,$18
		// ,$19
		// ,$20
		// ,$21
		// ,$22
		// ,$23
		// ,$24
		`) RETURNING id`,
		model.Username,
		model.Role,
		model.Password,
		// model.PasswordResetToken,
		// model.PasswordResetSentAt,
		model.Email,
		model.EmailUnconfirmed,
		model.EmailConfirmationToken,
		model.EmailConfirmationSentAt,
		// model.EmailConfirmedAt,
		// model.OTPEnabled,
		// model.OTPSecret,
		// model.OTPTimestep,
		// model.SignInLastAt,
		// model.SignInFailedAttempts,
		// model.SignInLockedAt,
		// model.SignInUnlockToken,
		// model.ProfilePicture,
		// model.Signature,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.ConfirmedAt,
		// model.BannedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
