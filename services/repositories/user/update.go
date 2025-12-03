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
			,otp_enabled = $11
			,otp_secret = $12
			,otp_timestep = $13
			,sign_in_last_at = $14
			,sign_in_failed_attempts = $15
			,sign_in_locked_at = $16
			,sign_in_unlock_token = $17
			,profile_picture = $18
			,signature = $19
			,created_at = $20
			,updated_at = $21
			,confirmed_at = $22
			,banned_at = $23
			,deleted_at = $24
		WHERE id = $25`,
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
