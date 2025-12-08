package user

import (
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) Update(model *user.User) (err error) {
	_, err = repo.db.Exec(`UPDATE users SET
		 username =                    $1
		,role =                        $2
		,member_of =                   $3
		,password =                    $4
		,password_reset_token =        $5
		,password_reset_sent_at =      $6
		,email =                       $7
		,email_unconfirmed =           $8
		,email_confirmation_token =    $9
		,email_confirmation_sent_at = $10
		,email_confirmed_at =         $11
		,language =                   $12
		,otp_enabled =                $13
		,otp_secret =                 $14
		,otp_timestep =               $15
		,sign_in_last_at =            $16
		,sign_in_failed_attempts =    $17
		,sign_in_locked_at =          $18
		,sign_in_unlock_token =       $19
		,profile_picture =            $20
		,signature =                  $21
		,created_at =                 $22
		,updated_at =                 $23
		,confirmed_at =               $24
		,banned_at =                  $25
		,deleted_at =                 $26
		WHERE id =                    $27`,
		model.Username,
		model.Role,
		model.MemberOf,
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
