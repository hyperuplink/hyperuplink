package activity

func (repo *Repository) DeleteOlderThanDays(
	kinds []string,
	days int,
) (deleted int64, err error) {
	if len(kinds) == 0 || days < 1 {
		return 0, nil
	}

	tag, err := repo.db.Exec(`DELETE FROM activities
		WHERE kind = ANY($1)
		AND created_at < NOW() - make_interval(days => $2)`,
		kinds,
		days,
	)
	if err != nil {
		return 0, repo.db.ConvertError(err)
	}

	return tag.RowsAffected(), nil
}
