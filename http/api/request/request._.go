package request

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logicperms "xn--gckvb8fzb.com/hyperuplink/logic/helpers/perms"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const UserLocal = "api_user"

type Request struct {
	r      route.IRouteController
	c      fiber.Ctx
	User   *user.User
	System *setting.System
	perms  *permission.Resolution
	lang   string
}

func New(
	r route.IRouteController,
	c fiber.Ctx,
) *Request {
	req := new(Request)
	req.r = r
	req.c = c

	if usr, ok := c.Locals(UserLocal).(*user.User); ok {
		req.User = usr
	}

	settingSystem, err := settingRepo.GetByID[setting.System](
		r.GetRuntime().Repositories.Setting,
		"system",
	)
	if err != nil {
		r.GetRuntime().Error("error", err)
		defaultSystem := setting.System{}
		req.System = &defaultSystem
	} else {
		req.System = &settingSystem.JSONValue
	}

	al := c.Get("Accept-Language", "en")
	als := strings.Split(al, "-")
	req.lang = strings.ToLower(als[0])

	return req
}

func (req *Request) Ts(msg string) string {
	return req.r.GetRuntime().Intnat.NewLocalizer(req.lang).Get(msg)
}

func (req *Request) Role() user.Role {
	if req.User == nil {
		return user.GuestRole
	}
	return req.User.Role
}

func (req *Request) UserUUID() (uuid.UUID, bool) {
	if req.User == nil {
		return uuid.Nil, false
	}
	return req.User.ID, true
}

func (req *Request) Perms() *permission.Resolution {
	if req.perms != nil {
		return req.perms
	}

	var memberOf []string
	if req.User != nil {
		memberOf = req.User.MemberOf
	}

	req.perms = logicperms.Resolve(
		req.r.GetRuntime(),
		req.Role(),
		memberOf,
	)
	return req.perms
}

func (req *Request) AccessControl(roles ...user.Role) (mustReturn bool, err error) {
	if slices.Index(roles, req.Role()) < 0 {
		return true, req.RespondError(errs.ErrForbidden)
	}

	return false, nil
}

func (req *Request) BindJSON(in any) (mustReturn bool, err error) {
	if berr := req.c.Bind().JSON(in); berr != nil {
		return true, req.respondBindError(berr, "json")
	}

	return false, nil
}

func (req *Request) BindForm(in any) (mustReturn bool, err error) {
	if berr := req.c.Bind().Form(in); berr != nil {
		return true, req.respondBindError(berr, "form")
	}

	return false, nil
}

func (req *Request) respondBindError(berr error, tag string) error {
	var valErrs validator.ValidationErrors
	if errors.As(berr, &valErrs) {
		fields := make(map[string]string)
		for _, e := range valErrs {
			fields[strings.ToLower(e.Field())] = fmt.Sprintf(
				"%s_%s_%s",
				errs.ErrValidation.Error(),
				strings.ToLower(e.Field()),
				e.Tag(),
			)
		}

		return req.c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":  errs.ErrValidation.Error(),
			"fields": fields,
		})
	}

	return req.c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": errs.ErrFormInvalid.Error(),
	})
}

func (req *Request) Respond(data any) error {
	return req.c.JSON(data)
}

func (req *Request) RespondCreated(data any) error {
	return req.c.Status(fiber.StatusCreated).JSON(data)
}

func (req *Request) RespondOK() error {
	return req.c.JSON(fiber.Map{"status": "ok"})
}

func (req *Request) RespondError(err error) error {
	return req.c.Status(statusForError(err)).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func (req *Request) RespondOnError(err error) (mustReturn bool, rerr error) {
	if err == nil {
		return false, nil
	}

	return true, req.RespondError(err)
}

// TODO: Find a different solution for this, e.g. integrating the fiber.Status
// into the error itself.
func statusForError(err error) int {
	switch {
	case errors.Is(err, errs.ErrNoRows),
		errors.Is(err, errs.ErrUserIDNotFound),
		errors.Is(err, errs.ErrTargetIDNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, errs.ErrUnauthorized),
		errors.Is(err, errs.ErrAPIKeyInvalid):
		return fiber.StatusUnauthorized
	case errors.Is(err, errs.ErrForbidden):
		return fiber.StatusForbidden
	case errors.Is(err, errs.ErrUniqueViolationOn):
		return fiber.StatusConflict
	case errors.Is(err, errs.ErrValidation),
		errors.Is(err, errs.ErrFormInvalid),
		errors.Is(err, errs.ErrPasswordWrong),
		errors.Is(err, errs.ErrOTPCodeWrong),
		errors.Is(err, errs.ErrOTPSetupExpired),
		errors.Is(err, errs.ErrInvalidUploadFormat),
		errors.Is(err, errs.ErrInvalidStorageProvider),
		errors.Is(err, errs.ErrInvalidTheme),
		errors.Is(err, errs.ErrInvalidColorscheme),
		errors.Is(err, errs.ErrThemeStorageNotConfigured),
		errors.Is(err, errs.ErrBackgroundUploadFailed),
		errors.Is(err, errs.ErrBackgroundAlreadySet),
		errors.Is(err, errs.ErrBackgroundFormatNotAllowed),
		errors.Is(err, errs.ErrBannerUploadFailed),
		errors.Is(err, errs.ErrBannerAlreadySet),
		errors.Is(err, errs.ErrBannerFormatNotAllowed),
		errors.Is(err, errs.ErrFaviconUploadFailed),
		errors.Is(err, errs.ErrFaviconAlreadySet),
		errors.Is(err, errs.ErrFaviconFormatNotAllowed),
		errors.Is(err, errs.ErrPictureTooLarge),
		errors.Is(err, errs.ErrPictureFormatNotAllowed),
		errors.Is(err, errs.ErrAttachmentTooLarge),
		errors.Is(err, errs.ErrAttachmentFormatNotAllowed),
		errors.Is(err, errs.ErrAttachmentHookFailed),
		errors.Is(err, errs.ErrAttachmentDuplicate),
		errors.Is(err, errs.ErrAttachmentNotFound),
		errors.Is(err, errs.ErrAttachmentUploadFailed),
		errors.Is(err, errs.ErrPollKindInvalid),
		errors.Is(err, errs.ErrPollEnded),
		errors.Is(err, errs.ErrPollSelectionInvalid),
		errors.Is(err, errs.ErrPollNotAllowed),
		errors.Is(err, errs.ErrPollOptionsTooFew),
		errors.Is(err, errs.ErrPollOptionsTooMany),
		errors.Is(err, errs.ErrPollOptionTooLong),
		errors.Is(err, errs.ErrPollEndsAtInvalid),
		errors.Is(err, errs.ErrPollEndsAtPast):
		return fiber.StatusUnprocessableEntity
	default:
		return fiber.StatusInternalServerError
	}
}
