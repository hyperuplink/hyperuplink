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
		,timezone =                   $13
		,otp_enabled =                $14
		,otp_secret =                 $15
		,otp_timestep =               $16
		,sign_in_last_at =            $17
		,sign_in_failed_attempts =    $18
		,sign_in_locked_at =          $19
		,sign_in_unlock_token =       $20
		,profile_picture =            $21
		,signature_text =             $22
		,signature_html =             $23
		,created_at =                 $24
		,updated_at =                 $25
		,confirmed_at =               $26
		,banned_at =                  $27
		,deleted_at =                 $28
		WHERE id =                    $29`,
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
		model.Timezone,
		model.OTPEnabled,
		model.OTPSecret,
		model.OTPTimestep,
		model.SignInLastAt,
		model.SignInFailedAttempts,
		model.SignInLockedAt,
		model.SignInUnlockToken,
		model.ProfilePicture,
		model.SignatureText,
		model.SignatureHTML,
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
