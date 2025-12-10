package user

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) Create(model *user.User) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO users (
		 username
		,role
		,member_of
		,password
		,email
		,email_unconfirmed
		,email_confirmation_token
		,email_confirmation_sent_at
		,language
		,created_at
		,updated_at
		) VALUES (
		 $1
		,$2
		,$3
		,$4
		,$5
		,$6
		,$7
		,$8
		,$9
		,NOW()
		,NOW()
		) RETURNING id`,
		model.Username,
		model.Role,
		model.MemberOf,
		model.Password,
		model.Email,
		model.EmailUnconfirmed,
		model.EmailConfirmationToken,
		model.EmailConfirmationSentAt,
		model.Language,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
