package postevent

import (
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
)

func (repo *Repository) Create(model *postevent.PostEvent) (rowID string, err error) {
	// TODO: Add correct attributes
	err = repo.db.QueryRow(`INSERT INTO postevents (
		 id
		,created_at
		,updated_at
		,deleted_at
		) VALUES (
		,$1
		,NOW()
		,NOW()
		,NULL
		) RETURNING id`,
		model.ID,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
