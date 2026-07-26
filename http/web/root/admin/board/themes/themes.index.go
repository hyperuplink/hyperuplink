package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicthemes.View
	view, err = logicthemes.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_theme", view.Theme)
	req.SetData("themes", view.Themes)
	req.SetData("colorschemes", view.Colorschemes)
	req.SetData("storage_ids", view.StorageIDs)

	return req.Respond()
}
