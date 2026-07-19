package apikey

import (
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) Delete(model *apikey.APIKey) (err error) {
	_, err = repo.db.Exec(
		`UPDATE apikeys SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`,
		model.ID,
	)
	return repo.db.ConvertError(err)
}
