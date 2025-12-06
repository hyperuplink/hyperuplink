package user

import (
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) Update(model *user.User) (err error) {
	_, err = repo.db.Exec(`UPDATE users SET
			 username = $1
			,role = $2
			,password = $3
			,password_reset_token = $4
			,password_reset_sent_at = $5
			,email = $6
			,email_unconfirmed = $7
			,email_confirmation_token = $8
			,email_confirmation_sent_at = $9
			,email_confirmed_at = $10
			,language = $11
			,otp_enabled = $12
			,otp_secret = $13
			,otp_timestep = $14
			,sign_in_last_at = $15
			,sign_in_failed_attempts = $16
			,sign_in_locked_at = $17
			,sign_in_unlock_token = $18
			,profile_picture = $19
			,signature = $20
			,created_at = $21
			,updated_at = $22
			,confirmed_at = $23
			,banned_at = $24
			,deleted_at = $25
		WHERE id = $26`,
		model.Username,
		model.Role,
		model.Password,
		model.PasswordResetToken,
		model.PasswordResetSentAt,
		model.Email,
		model.EmailUnconfirmed,
		model.EmailConfirmationToken,
		model.EmailConfirmationSentAt,
		model.EmailConfirmedAt,
		model.Language,
		model.OTPEnabled,
		model.OTPSecret,
		model.OTPTimestep,
		model.SignInLastAt,
		model.SignInFailedAttempts,
		model.SignInLockedAt,
		model.SignInUnlockToken,
		model.ProfilePicture,
		model.Signature,
		model.CreatedAt,
		model.UpdatedAt,
		model.ConfirmedAt,
		model.BannedAt,
		model.DeletedAt,
		// WHERE id =
		model.ID,
	)
	return repo.db.ConvertError(err)
}
