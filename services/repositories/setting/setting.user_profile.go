package setting

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func GetOrCreateUserProfile(
	ctx context.Context,
	repo *Repository,
	userID uuid.UUID,
) (model *setting.Setting[setting.UserProfile], err error) {
	id := setting.UserProfileID(userID)

	model, err = GetByID[setting.UserProfile](ctx, repo, id)
	if err == nil {
		return model, nil
	}
	if !errors.Is(err, errs.ErrNoRows) {
		return nil, err
	}

	model = new(setting.Setting[setting.UserProfile])
	model.ID = id
	model.JSONValue = setting.NewUserProfile()

	if _, err = Create(ctx, repo, model); err != nil {
		return nil, err
	}

	return model, nil
}
