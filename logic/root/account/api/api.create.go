package api

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

type CreateInput struct {
	KeyName string `json:"name" form:"name" validate:"required,min=1,max=64"`
}

func Create(
	rt *runtime.Runtime,
	userID uuid.UUID,
	in *CreateInput,
) (key *apikey.APIKey, secret string, err error) {
	if key, secret, err = apikey.New(userID, in.KeyName); err != nil {
		return nil, "", err
	}

	if key.ID, err = rt.Repositories.APIKey.Create(key); err != nil {
		return nil, "", err
	}

	return key, secret, nil
}
