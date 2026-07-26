package general

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	Name                  string `json:"name" form:"name" validate:"required,min=1,max=128"`
	BaseURL               string `json:"base_url" form:"base_url" validate:"required,url,max=255"`
	TopicsPerPage         int    `json:"topics_per_page" form:"topics_per_page" validate:"required,min=1,max=200"`
	PostsPerPage          int    `json:"posts_per_page" form:"posts_per_page" validate:"required,min=1,max=200"`
	AdminLogRetentionDays int    `json:"admin_log_retention_days" form:"admin_log_retention_days" validate:"required,min=1,max=3650"`

	EnableAbout bool   `json:"enable_about" form:"enable_about"`
	About       string `json:"about" form:"about" validate:"required_if=EnableAbout true"`

	EnableContact bool   `json:"enable_contact" form:"enable_contact"`
	Contact       string `json:"contact" form:"contact" validate:"required_if=EnableContact true"`

	EnablePrivacyPolicy bool   `json:"enable_privacy_policy" form:"enable_privacy_policy"`
	PrivacyPolicy       string `json:"privacy_policy" form:"privacy_policy" validate:"required_if=EnablePrivacyPolicy true"`

	EnableTerms bool   `json:"enable_terms" form:"enable_terms"`
	Terms       string `json:"terms" form:"terms" validate:"required_if=EnableTerms true"`

	EnableQuit bool   `json:"enable_quit" form:"enable_quit"`
	QuitURL    string `json:"quit_url" form:"quit_url" validate:"required_if=EnableQuit true,omitempty,url,max=255"`
}

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	settingSystem, err := settingRepo.GetByID[setting.System](
		gh.Repositories(rt).Setting,
		"system",
	)
	if err != nil {
		return err
	}

	beforeSystem := settingSystem.JSONValue

	settingSystem.JSONValue.Name = in.Name
	settingSystem.JSONValue.BaseURL = in.BaseURL
	settingSystem.JSONValue.TopicsPerPage = in.TopicsPerPage
	settingSystem.JSONValue.PostsPerPage = in.PostsPerPage
	settingSystem.JSONValue.AdminLogRetentionDays = in.AdminLogRetentionDays

	if err = settingRepo.Update[setting.System](
		gh.Repositories(rt).Setting,
		settingSystem,
	); err != nil {
		return err
	}

	settingGeneral, err := settingRepo.GetByID[setting.General](
		gh.Repositories(rt).Setting,
		"general",
	)
	if err != nil {
		return err
	}

	beforeGeneral := settingGeneral.JSONValue

	settingGeneral.JSONValue.EnableAbout = in.EnableAbout
	settingGeneral.JSONValue.About = in.About
	settingGeneral.JSONValue.EnableContact = in.EnableContact
	settingGeneral.JSONValue.Contact = in.Contact
	settingGeneral.JSONValue.EnablePrivacyPolicy = in.EnablePrivacyPolicy
	settingGeneral.JSONValue.PrivacyPolicy = in.PrivacyPolicy
	settingGeneral.JSONValue.EnableTerms = in.EnableTerms
	settingGeneral.JSONValue.Terms = in.Terms
	settingGeneral.JSONValue.EnableQuit = in.EnableQuit
	settingGeneral.JSONValue.QuitURL = in.QuitURL

	if err = settingRepo.Update[setting.General](
		gh.Repositories(rt).Setting,
		settingGeneral,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"system", beforeSystem, settingSystem.JSONValue)
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"general", beforeGeneral, settingGeneral.JSONValue)
	}

	return nil
}
