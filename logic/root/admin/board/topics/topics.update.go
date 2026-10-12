package topics

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	AllowKindQuestion bool `json:"allow_kind_question" form:"allow_kind_question"`
	AllowKindPoll     bool `json:"allow_kind_poll" form:"allow_kind_poll"`
	AllowKindRSVP     bool `json:"allow_kind_rsvp" form:"allow_kind_rsvp"`
}

func Update(
	ctx context.Context,
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	settingTopics, err := settingRepo.GetByID[setting.Topics](ctx,
		gh.Repositories(rt).Setting,
		"topics",
	)
	if err != nil {
		return err
	}

	before := settingTopics.JSONValue

	settingTopics.JSONValue.AllowKindQuestion = in.AllowKindQuestion
	settingTopics.JSONValue.AllowKindPoll = in.AllowKindPoll
	settingTopics.JSONValue.AllowKindRSVP = in.AllowKindRSVP

	if err = settingRepo.Update[setting.Topics](ctx,
		gh.Repositories(rt).Setting,
		settingTopics,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(ctx, rt, actorID.UUID,
			"topics", before, settingTopics.JSONValue)
	}

	return nil
}
