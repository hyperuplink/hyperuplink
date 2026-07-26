package postevent

import (
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vpostevent"
)

func (repo *Repository) VAllReports(
	qo common.QueryOptions,
) (model *[]vpostevent.VPostEvent, err error) {
	var rows pgx.Rows
	var mod []vpostevent.VPostEvent

	rows, err = repo.db.Query(qo.Query(
		`SELECT
			 pe.*
			,COALESCE(t.text, r.text, '') AS target_text
			,COALESCE(tu.username, ru.username, '') AS target_author_username
			,COALESCE(au.username, '') AS author_username
			,COALESCE(c.slug, '') AS category_slug
			,COALESCE(f.slug, '') AS forum_slug
			,COALESCE(ot.slug, '') AS topic_slug
			,COALESCE(t.short_id, r.short_id, '') AS target_short_id
		FROM postevents pe
		LEFT JOIN topics t ON pe.target = 'topic' AND t.id = pe.topic_id
		LEFT JOIN replies r ON pe.target = 'reply' AND r.id = pe.reply_id
		LEFT JOIN users tu ON tu.id = t.author_id
		LEFT JOIN users ru ON ru.id = r.author_id
		LEFT JOIN users au ON au.id = pe.author_id
		LEFT JOIN topics ot ON ot.id = COALESCE(t.id, r.topic_id)
		LEFT JOIN forums f ON f.id = ot.forum_id
		LEFT JOIN categories c ON c.id = f.category_id
		WHERE pe.type = 'report'`,
		common.QueryCapabilities{
			Table:      "pe",
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vpostevent.VPostEvent])

	return &mod, repo.db.ConvertError(err)
}
