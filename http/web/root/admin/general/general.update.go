package general

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

type GeneralUpdateForm struct {
	Name          string `form:"name" validate:"required,min=1,max=128"`
	BaseURL       string `form:"base_url" validate:"required,url,max=255"`
	TopicsPerPage int    `form:"topics_per_page" validate:"required,min=1,max=200"`
	PostsPerPage  int    `form:"posts_per_page" validate:"required,min=1,max=200"`

	EnableAbout bool   `form:"enable_about"`
	About       string `form:"about" validate:"required_if=EnableAbout true"`

	EnableContact bool   `form:"enable_contact"`
	Contact       string `form:"contact" validate:"required_if=EnableContact true"`

	EnablePrivacyPolicy bool   `form:"enable_privacy_policy"`
	PrivacyPolicy       string `form:"privacy_policy" validate:"required_if=EnablePrivacyPolicy true"`

	EnableTerms bool   `form:"enable_terms"`
	Terms       string `form:"terms" validate:"required_if=EnableTerms true"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminGeneral")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(GeneralUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	var settingSystem *setting.Setting[setting.System]
	settingSystem, err = settingRepo.GetByID[setting.System](
		r.Runtime.Repositories.Setting,
		"system",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingSystem.JSONValue.Name = frm.Name
	settingSystem.JSONValue.BaseURL = frm.BaseURL
	settingSystem.JSONValue.TopicsPerPage = frm.TopicsPerPage
	settingSystem.JSONValue.PostsPerPage = frm.PostsPerPage

	err = settingRepo.Update[setting.System](
		r.Runtime.Repositories.Setting,
		settingSystem,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var settingGeneral *setting.Setting[setting.General]
	settingGeneral, err = settingRepo.GetByID[setting.General](
		r.Runtime.Repositories.Setting,
		"general",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingGeneral.JSONValue.EnableAbout = frm.EnableAbout
	settingGeneral.JSONValue.About = frm.About
	settingGeneral.JSONValue.EnableContact = frm.EnableContact
	settingGeneral.JSONValue.Contact = frm.Contact
	settingGeneral.JSONValue.EnablePrivacyPolicy = frm.EnablePrivacyPolicy
	settingGeneral.JSONValue.PrivacyPolicy = frm.PrivacyPolicy
	settingGeneral.JSONValue.EnableTerms = frm.EnableTerms
	settingGeneral.JSONValue.Terms = frm.Terms

	err = settingRepo.Update[setting.General](
		r.Runtime.Repositories.Setting,
		settingGeneral,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
