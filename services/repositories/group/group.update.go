package group

import (
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Update(model *group.Group) (err error) {
	_, err = repo.db.Exec(`UPDATE groups SET
		 name =                        $1
		,updated_at =                  NOW()
		WHERE id =                     $2`,
		model.Name,
		// WHERE id =
		model.ID,
	)
	return repo.db.ConvertError(err)
}
