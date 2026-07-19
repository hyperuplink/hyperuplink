package apikey

import (
	"github.com/google/uuid"
)

func (repo *Repository) TouchLastUsed(id uuid.UUID) (err error) {
	_, err = repo.db.Exec(
		`UPDATE apikeys SET last_used_at = NOW() WHERE id = $1`,
		id,
	)
	return repo.db.ConvertError(err)
}
