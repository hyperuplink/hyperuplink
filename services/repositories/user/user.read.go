package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]user.User, err error) {
	var rows pgx.Rows
	var mod []user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE id = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}

func (repo *Repository) GetByUsername(
	username string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE username = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		username)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByEmail(
	email string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE email = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		email)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllToNotifyForReply(
	topicID uuid.UUID,
	replyID uuid.UUID,
	authorID uuid.UUID,
	categoryID uuid.UUID,
	qo common.QueryOptions,
) (model *[]user.User, err error) {
	var rows pgx.Rows
	var mod []user.User

	rows, err = repo.db.Query(qo.Query(
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
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}
