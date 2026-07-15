package setting

import (
	"errors"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func GetOrCreateUserProfile(
	repo *Repository,
	userID uuid.UUID,
) (model *setting.Setting[setting.UserProfile], err error) {
	id := setting.UserProfileID(userID)

	model, err = GetByID[setting.UserProfile](repo, id)
	if err == nil {
		return model, nil
	}
	if !errors.Is(err, errs.ErrNoRows) {
		return nil, err
	}

	model = new(setting.Setting[setting.UserProfile])
	model.ID = id
	model.JSONValue = setting.NewUserProfile()

	if _, err = Create(repo, model); err != nil {
		return nil, err
	}

	return model, nil
}
