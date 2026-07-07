package search

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/vsearchresult"
)

type Options struct {
	Title       bool   // topic title (topics.name)
	Body        bool   // topic body (topics.text)
	Replies     bool   // reply bodies (replies.text)
	Attachments bool   // placeholder, attachments not searchable yet
	Author      string // when non-empty, restrict to this username
	// AllowedCategorySlugs restricts results to these category slugs. A nil
	// value means "no restriction" (may read every category), a non-nil but
	// empty slice means "may read nothing" and yields zero results.
	AllowedCategorySlugs []string
}

const topicColumns = `
	'topic' AS kind,
	vt.category_slug,
	vt.forum_slug,
	vt.slug AS topic_slug,
	vt.name AS topic_name,
	vt.author_username,
	vt.created_at,
	'' AS reply_short_id,
	0::bigint AS reply_position,
	vt.updated_at AS sort_at`

const replyColumns = `
	'reply' AS kind,
	c.slug AS category_slug,
	f.slug AS forum_slug,
	t.slug AS topic_slug,
	t.name AS topic_name,
	vr.author_username,
	vr.created_at,
	vr.short_id AS reply_short_id,
	-- 1-based position of this reply among the topic's visible replies,
	-- ordered by created_at (id as a stable tie-break).
	(SELECT COUNT(*) FROM replies r2
	   WHERE r2.topic_id = vr.topic_id
	     AND r2.deleted_at IS NULL AND r2.spammed_at IS NULL
	     AND (r2.created_at < vr.created_at
	          OR (r2.created_at = vr.created_at AND r2.id <= vr.id))
	) AS reply_position,
	vr.created_at AS sort_at`

const topicFrom = ` FROM vtopics vt`

const replyFrom = ` FROM vreplies vr
	JOIN topics t ON t.id = vr.topic_id
	JOIN forums f ON f.id = t.forum_id
	JOIN categories c ON c.id = f.category_id`

const topicBaseWhere = `vt.spammed_at IS NULL AND vt.deleted_at IS NULL`

const replyBaseWhere = `vr.spammed_at IS NULL AND vr.deleted_at IS NULL
	AND t.deleted_at IS NULL AND t.spammed_at IS NULL
	AND f.deleted_at IS NULL AND c.deleted_at IS NULL`

func buildArgs(term string, opts Options) (args []any, authorPh string, catPh string) {
	args = []any{term}
	if opts.Author != "" {
		args = append(args, opts.Author)
		authorPh = fmt.Sprintf("$%d", len(args))
	}
	if opts.AllowedCategorySlugs != nil {
		args = append(args, opts.AllowedCategorySlugs)
		catPh = fmt.Sprintf("$%d", len(args))
	}
	return args, authorPh, catPh
}

func branchWheres(opts Options, authorPh string, catPh string) (topicWhere string, hasTopic bool, replyWhere string, hasReply bool) {
	var cols []string
	if opts.Title {
		cols = append(cols, "coalesce(vt.name, '')")
	}
	if opts.Body {
		cols = append(cols, "coalesce(vt.text, '')")
	}
	if len(cols) > 0 {
		fts := fmt.Sprintf(
			"to_tsvector('english', %s) @@ websearch_to_tsquery('english', $1)",
			strings.Join(cols, " || ' ' || "),
		)
		topicWhere = fts + " AND " + topicBaseWhere
		if authorPh != "" {
			topicWhere += " AND lower(vt.author_username) = lower(" + authorPh + ")"
		}
		if catPh != "" {
			topicWhere += " AND vt.category_slug = ANY(" + catPh + ")"
		}
		hasTopic = true
	}

	if opts.Replies {
		fts := "to_tsvector('english', coalesce(vr.text, '')) @@ websearch_to_tsquery('english', $1)"
		replyWhere = fts + " AND " + replyBaseWhere
		if authorPh != "" {
			replyWhere += " AND lower(vr.author_username) = lower(" + authorPh + ")"
		}
		if catPh != "" {
			replyWhere += " AND c.slug = ANY(" + catPh + ")"
		}
		hasReply = true
	}

	return topicWhere, hasTopic, replyWhere, hasReply
}

func (repo *Repository) Count(term string, opts Options) (total int64, err error) {
	args, authorPh, catPh := buildArgs(term, opts)
	topicWhere, hasTopic, replyWhere, hasReply := branchWheres(opts, authorPh, catPh)
	if !hasTopic && !hasReply {
		return 0, nil
	}

	var parts []string
	if hasTopic {
		parts = append(parts, "SELECT 1"+topicFrom+" WHERE "+topicWhere)
	}
	if hasReply {
		parts = append(parts, "SELECT 1"+replyFrom+" WHERE "+replyWhere)
	}

	var rows pgx.Rows
	rows, err = repo.db.Query(
		"SELECT COUNT(*) AS total FROM ("+strings.Join(parts, " UNION ALL ")+") results",
		args...,
	)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	pag, err = pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}
	total = pag["total"].(int64)

	return total, nil
}

func (repo *Repository) Query(
	term string,
	opts Options,
	limit int,
	page int,
) (model *[]vsearchresult.VSearchResult, total int64, err error) {
	total, err = repo.Count(term, opts)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}
	if total == 0 {
		empty := []vsearchresult.VSearchResult{}
		return &empty, 0, nil
	}

	args, authorPh, catPh := buildArgs(term, opts)
	topicWhere, hasTopic, replyWhere, hasReply := branchWheres(opts, authorPh, catPh)

	var parts []string
	if hasTopic {
		parts = append(parts, "SELECT"+topicColumns+topicFrom+" WHERE "+topicWhere)
	}
	if hasReply {
		parts = append(parts, "SELECT"+replyColumns+replyFrom+" WHERE "+replyWhere)
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	limitPh := fmt.Sprintf("$%d", len(args)+1)
	offsetPh := fmt.Sprintf("$%d", len(args)+2)
	args = append(args, limit, offset)

	var rows pgx.Rows
	rows, err = repo.db.Query(
		`SELECT kind, category_slug, forum_slug, topic_slug, topic_name,
		        author_username, created_at, reply_short_id, reply_position
		 FROM (`+strings.Join(parts, " UNION ALL ")+`) results
		 ORDER BY sort_at DESC
		 LIMIT `+limitPh+` OFFSET `+offsetPh,
		args...,
	)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	mod, err := pgx.CollectRows(rows, pgx.RowToStructByName[vsearchresult.VSearchResult])

	return &mod, total, repo.db.ConvertError(err)
}
