package attachments

import (
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type AttachmentsUpdateForm struct {
	EnableAttachments  bool     `form:"enable_attachments"`
	UploadFormats      []string `form:"upload_formats" validate:"required_if=EnableAttachments true"`
	MaxSize            int64    `form:"max_size" validate:"required,min=1"`
	StorageProviderID  string   `form:"storage_provider_id" validate:"required_if=EnableAttachments true,max=64"`
	StoragePath        string   `form:"storage_path" validate:"omitempty,max=255"`
	OnUploadHook       string   `form:"on_upload_hook" validate:"omitempty,max=1024"`
	InlineImageDisplay bool     `form:"inline_image_display"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardAttachments")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(AttachmentsUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	for _, format := range frm.UploadFormats {
		if !slices.Contains(setting.AttachmentUploadFormatOptions, format) {
			req.Flash.SetError(errs.ErrInvalidUploadFormat)
			return req.RedirectToRoute(myRoute)
		}
	}

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	valid := frm.StorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == frm.StorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		req.Flash.SetError(errs.ErrInvalidStorageProvider)
		return req.RedirectToRoute(myRoute)
	}

	var settingAttachments *setting.Setting[setting.Attachments]
	settingAttachments, err = settingRepo.GetByID[setting.Attachments](
		r.Runtime.Repositories.Setting,
		"attachments",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	before := settingAttachments.JSONValue

	settingAttachments.JSONValue.EnableAttachments = frm.EnableAttachments
	settingAttachments.JSONValue.UploadFormats = frm.UploadFormats
	settingAttachments.JSONValue.MaxSize = frm.MaxSize
	settingAttachments.JSONValue.StorageProviderID = frm.StorageProviderID
	settingAttachments.JSONValue.StoragePath = frm.StoragePath
	settingAttachments.JSONValue.OnUploadHook = frm.OnUploadHook
	settingAttachments.JSONValue.InlineImageDisplay = frm.InlineImageDisplay

	err = settingRepo.Update[setting.Attachments](
		r.Runtime.Repositories.Setting,
		settingAttachments,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if actorID, ok := req.Session.GetUserUUID(); ok {
		logicactivity.RecordAdminSettingsUpdate(r.Runtime, actorID,
			"attachments", before, settingAttachments.JSONValue)
	}

	return req.RedirectToRoute(myRoute)
}
