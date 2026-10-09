package problem

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
)

const (
	MediaType = "application/problem+json"
	TypeBlank = "about:blank"
	Realm     = "hyperuplink"

	retryAfterUnavailable = "5"
)

type Problem struct {
	Type     string              `json:"type"`
	Title    string              `json:"title"`
	Status   int                 `json:"status"`
	Detail   string              `json:"detail,omitempty"`
	Instance string              `json:"instance,omitempty"`
	Code     string              `json:"code"`
	Errors   map[string][]string `json:"errors,omitempty"`
}

type mapping struct {
	err    error
	status int
}

var mappings = []mapping{
	{errs.ErrNoRows, fiber.StatusNotFound},
	{errs.ErrUserIDNotFound, fiber.StatusNotFound},
	{errs.ErrTargetIDNotFound, fiber.StatusNotFound},
	{errs.ErrAttachmentNotFound, fiber.StatusNotFound},
	{errs.ErrUnauthorized, fiber.StatusUnauthorized},
	{errs.ErrAPIKeyInvalid, fiber.StatusUnauthorized},
	{errs.ErrAuthenticationRequired, fiber.StatusUnauthorized},
	{errs.ErrOTPRequired, fiber.StatusUnauthorized},
	{errs.ErrUsernamePasswordWrong, fiber.StatusUnauthorized},
	{errs.ErrForbidden, fiber.StatusForbidden},
	{errs.ErrUniqueViolationOn, fiber.StatusConflict},
	{errs.ErrRateLimited, fiber.StatusTooManyRequests},
	{errs.ErrUnavailable, fiber.StatusServiceUnavailable},
	{errs.ErrFormInvalid, fiber.StatusBadRequest},
	{errs.ErrValidation, fiber.StatusUnprocessableEntity},
	{errs.ErrPasswordWrong, fiber.StatusUnprocessableEntity},
	{errs.ErrOTPCodeWrong, fiber.StatusUnprocessableEntity},
	{errs.ErrOTPSetupExpired, fiber.StatusUnprocessableEntity},
	{errs.ErrOTPTooManyAttempts, fiber.StatusUnprocessableEntity},
	{errs.ErrInvalidUploadFormat, fiber.StatusUnprocessableEntity},
	{errs.ErrInvalidStorageProvider, fiber.StatusUnprocessableEntity},
	{errs.ErrInvalidTheme, fiber.StatusUnprocessableEntity},
	{errs.ErrInvalidColorscheme, fiber.StatusUnprocessableEntity},
	{errs.ErrThemeStorageNotConfigured, fiber.StatusUnprocessableEntity},
	{errs.ErrBackgroundUploadFailed, fiber.StatusUnprocessableEntity},
	{errs.ErrBackgroundAlreadySet, fiber.StatusUnprocessableEntity},
	{errs.ErrBackgroundFormatNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrBannerUploadFailed, fiber.StatusUnprocessableEntity},
	{errs.ErrBannerAlreadySet, fiber.StatusUnprocessableEntity},
	{errs.ErrBannerFormatNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrFaviconUploadFailed, fiber.StatusUnprocessableEntity},
	{errs.ErrFaviconAlreadySet, fiber.StatusUnprocessableEntity},
	{errs.ErrFaviconFormatNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrPictureTooLarge, fiber.StatusUnprocessableEntity},
	{errs.ErrPictureFormatNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentTooLarge, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentFormatNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentHookFailed, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentDuplicate, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentUploadFailed, fiber.StatusUnprocessableEntity},
	{errs.ErrPollKindInvalid, fiber.StatusUnprocessableEntity},
	{errs.ErrPollEnded, fiber.StatusUnprocessableEntity},
	{errs.ErrPollSelectionInvalid, fiber.StatusUnprocessableEntity},
	{errs.ErrPollNotAllowed, fiber.StatusUnprocessableEntity},
	{errs.ErrPollOptionsTooFew, fiber.StatusUnprocessableEntity},
	{errs.ErrPollOptionsTooMany, fiber.StatusUnprocessableEntity},
	{errs.ErrPollOptionTooLong, fiber.StatusUnprocessableEntity},
	{errs.ErrPollEndsAtInvalid, fiber.StatusUnprocessableEntity},
	{errs.ErrPollEndsAtPast, fiber.StatusUnprocessableEntity},
	{errs.ErrEmailConfirmationTokenWrong, fiber.StatusUnprocessableEntity},
	{errs.ErrAttachmentsDisabled, fiber.StatusUnprocessableEntity},
	{errs.ErrReplyParentInvalid, fiber.StatusUnprocessableEntity},
	{errs.ErrPicturesDisabled, fiber.StatusUnprocessableEntity},
}

func New(status int, code string) Problem {
	return Problem{
		Type:   TypeBlank,
		Title:  http.StatusText(status),
		Status: status,
		Code:   code,
	}
}

func From(err error) Problem {
	for _, m := range mappings {
		if !errors.Is(err, m.err) {
			continue
		}
		if m.err == errs.ErrValidation {
			return New(m.status, err.Error())
		}
		return New(m.status, m.err.Error())
	}

	return New(fiber.StatusInternalServerError, errs.ErrInternal.Error())
}

func WithStatus(status int, err error) Problem {
	return New(status, err.Error())
}

func Validation(fields map[string][]string) Problem {
	p := New(fiber.StatusUnprocessableEntity, errs.ErrValidation.Error())
	p.Errors = fields

	return p
}

func FromFiber(err error) Problem {
	var ferr *fiber.Error
	if errors.As(err, &ferr) {
		return New(ferr.Code, CodeForStatus(ferr.Code))
	}

	return From(err)
}

func CodeForStatus(status int) string {
	text := strings.ToLower(http.StatusText(status))
	if text == "" {
		return errs.ErrInternal.Error()
	}

	var b strings.Builder
	b.WriteString("err_")
	for _, r := range text {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-':
			b.WriteByte('_')
		}
	}

	return b.String()
}

func Challenge(p Problem) string {
	if p.Code == errs.ErrAPIKeyInvalid.Error() {
		return `Bearer realm="` + Realm + `", error="invalid_token"`
	}

	return `Bearer realm="` + Realm + `"`
}

func Respond(c fiber.Ctx, p Problem) error {
	p.Instance = c.Path()

	switch p.Status {
	case fiber.StatusUnauthorized:
		c.Set(fiber.HeaderWWWAuthenticate, Challenge(p))
	case fiber.StatusServiceUnavailable:
		c.Set(fiber.HeaderRetryAfter, retryAfterUnavailable)
	}

	return c.Status(p.Status).JSON(p, MediaType)
}
