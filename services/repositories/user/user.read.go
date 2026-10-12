package user

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (repo *Repository) All(
	ctx context.Context,
	qo common.QueryOptions,
) (model []user.User, err error) {
	return database.Rows[user.User](ctx, repo.db, qo.Query(
		`SELECT * FROM users`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
	)
}

func (repo *Repository) GetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *user.User, err error) {
	return database.Row[user.User](ctx, repo.db, qo.Query(
		`SELECT * FROM users WHERE id = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) GetByID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, database.ConvertError(err)
	}

	return repo.GetByUUID(ctx, uuID, qo)
}

func (repo *Repository) GetByUsername(
	ctx context.Context,
	username string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	return database.Row[user.User](ctx, repo.db, qo.Query(
		`SELECT * FROM users WHERE username = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		username,
	)
}

func (repo *Repository) GetByEmail(
	ctx context.Context,
	email string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	return database.Row[user.User](ctx, repo.db, qo.Query(
		`SELECT * FROM users WHERE email = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		email,
	)
}

func (repo *Repository) AllToNotifyForReply(
	ctx context.Context,
	topicID uuid.UUID,
	replyID uuid.UUID,
	authorID uuid.UUID,
	categoryID uuid.UUID,
	qo common.QueryOptions,
) (model []user.User, err error) {
	return database.Rows[user.User](ctx, repo.db, qo.Query(
		`SELECT u.* FROM users u
		JOIN (
			SELECT author_id FROM topics WHERE id = $1
			UNION
			SELECT author_id FROM replies
			WHERE topic_id = $1 AND id <> $2 AND deleted_at IS NULL
		) p ON p.author_id = u.id
		LEFT JOIN settings s ON s.id = $7 || u.id::text
		WHERE u.id <> $3
		AND COALESCE((s.json_value ->> $8)::boolean, $9)
		AND (
			u.role = $4
			OR GREATEST(
				COALESCE((SELECT bits::int FROM permissions
					WHERE group_id IS NULL AND category_id = $5
					AND deleted_at IS NULL),
				(SELECT bits::int FROM permissions
					WHERE group_id IS NULL AND category_id IS NULL
					AND deleted_at IS NULL), 0),
				COALESCE((SELECT MAX(bits::int) FROM permissions
					WHERE group_id = ANY(u.member_of) AND category_id = $5
					AND deleted_at IS NULL), 0)
			) >= $6
		)`,
		common.QueryCapabilities{
			Table:      "u",
			HasBanned:  true,
			HasDeleted: true,
		}),
		topicID,
		replyID,
		authorID,
		string(user.AdminRole),
		categoryID,
		int(permission.ReadOnly),
		setting.UserProfilePrefix,
		setting.UserProfileNotifyOnReplyKey,
		setting.NewUserProfile().NotifyOnReply,
	)
}

func (repo *Repository) CountAll(ctx context.Context, qo common.QueryOptions) (total int64, err error) {

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	return database.Value[int64](ctx, repo.db, qoc.Query(
		`SELECT COUNT(id) AS total FROM users`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
	)
}
