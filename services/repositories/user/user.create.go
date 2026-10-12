package user

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (repo *Repository) Create(ctx context.Context, model *user.User) (rowID uuid.UUID, err error) {
	return database.Value[uuid.UUID](ctx, repo.db, `INSERT INTO users (
		 username
		,role
		,member_of
		,password
		,email
		,email_unconfirmed
		,email_confirmation_token
		,email_confirmation_sent_at
		,email_is_jid
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
		,$10
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
		model.EmailIsJID,
		model.Language,
	)
}
