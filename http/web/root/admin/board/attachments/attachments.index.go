package attachments

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/attachments"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardAttachments")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicattachments.View
	view, err = logicattachments.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_attachments", view.Attachments)
	req.SetData("storage_ids", view.StorageIDs)
	req.SetData("storage_is_public", view.StorageIsPublic)
	req.SetData("upload_format_options", view.UploadFormatOptions)

	return req.Respond()
}
