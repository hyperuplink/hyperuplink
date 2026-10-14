package errs

import (
	"errors"

	glideserrs "xn--gckvb8fzb.com/glides/errs"
)

var (
	ErrNotImplemented          = glideserrs.ErrNotImplemented
	ErrConfigTypeUnsupported   = glideserrs.ErrConfigTypeUnsupported
	ErrNoRows                  = glideserrs.ErrNoRows
	ErrUniqueViolationOn       = glideserrs.ErrUniqueViolationOn
	ErrUnavailable             = glideserrs.ErrUnavailable
	ErrHashInvalid             = glideserrs.ErrHashInvalid
	ErrHashVariantIncompatible = glideserrs.ErrHashVariantIncompatible
	ErrHashVersionIncompatible = glideserrs.ErrHashVersionIncompatible
	ErrTargetIDNotFound        = glideserrs.ErrTargetIDNotFound
	ErrNoSuchTargetType        = glideserrs.ErrNoSuchTargetType
	ErrStorageIDNotFound       = glideserrs.ErrStorageIDNotFound
	ErrStorageTypeInvalid      = glideserrs.ErrStorageTypeInvalid
	ErrFilePathInvalid         = glideserrs.ErrFilePathInvalid
	ErrJobTypeInvalid          = glideserrs.ErrJobTypeInvalid
	ErrJobSubTypeInvalid       = glideserrs.ErrJobSubTypeInvalid
	ErrJobPayloadInvalid       = glideserrs.ErrJobPayloadInvalid
	ErrCronFunctionIDExists    = glideserrs.ErrCronFunctionIDExists
	ErrCronFunctionIDNotFound  = glideserrs.ErrCronFunctionIDNotFound
	ErrCronFunctionInvalid     = glideserrs.ErrCronFunctionInvalid
)

var (
	ErrUserIDNotFound error = newPublic(
		"err_user_id_not_found",
	)
	ErrPasswordWrong error = newPublic(
		"err_password_wrong",
	)
	ErrIfaceTypeUnsupported error = errors.New(
		"err_iface_type_unsupported",
	)
	ErrRedisAddrsEmpty error = errors.New(
		"err_redis_addrs_empty",
	)
	ErrRedisAddrsMalformed error = errors.New(
		"err_redis_addrs_malformed",
	)
	ErrFormInvalid error = newPublic(
		"err_form_invalid",
	)
	ErrUsernamePasswordWrong error = newPublic(
		"username_password_wrong",
	)
	ErrEmailConfirmationTokenWrong error = newPublic(
		"email_confirmation_token_wrong",
	)
	ErrPictureTooLarge error = newPublic(
		"picture_too_large",
	)
	ErrPictureFormatNotAllowed error = newPublic(
		"picture_format_not_allowed",
	)
	ErrInvalidUploadFormat error = newPublic(
		"invalid_upload_format",
	)
	ErrInvalidStorageProvider error = newPublic(
		"invalid_storage_provider",
	)
	ErrThemeStorageNotConfigured error = newPublic(
		"theme_storage_not_configured",
	)
	ErrInvalidTheme error = newPublic(
		"invalid_theme",
	)
	ErrInvalidColorscheme error = newPublic(
		"invalid_colorscheme",
	)
	ErrBannerUploadFailed error = newPublic(
		"banner_upload_failed",
	)
	ErrBannerAlreadySet error = newPublic(
		"banner_already_set",
	)
	ErrBannerFormatNotAllowed error = newPublic(
		"banner_format_not_allowed",
	)
	ErrFaviconUploadFailed error = newPublic(
		"favicon_upload_failed",
	)
	ErrFaviconAlreadySet error = newPublic(
		"favicon_already_set",
	)
	ErrFaviconFormatNotAllowed error = newPublic(
		"favicon_format_not_allowed",
	)
	ErrBackgroundUploadFailed error = newPublic(
		"background_upload_failed",
	)
	ErrBackgroundAlreadySet error = newPublic(
		"background_already_set",
	)
	ErrBackgroundFormatNotAllowed error = newPublic(
		"background_format_not_allowed",
	)
	ErrAttachmentTooLarge error = newPublic(
		"attachment_too_large",
	)
	ErrAttachmentUploadFailed error = newPublic(
		"attachment_upload_failed",
	)
	ErrAttachmentFormatNotAllowed error = newPublic(
		"attachment_format_not_allowed",
	)
	ErrAttachmentHookFailed error = newPublic(
		"attachment_hook_failed",
	)
	ErrAttachmentDuplicate error = newPublic(
		"attachment_duplicate",
	)
	ErrAttachmentNotFound error = newPublic(
		"attachment_not_found",
	)
	ErrAPIKeyInvalid error = newPublic(
		"err_apikey_invalid",
	)
	ErrForbidden error = newPublic(
		"err_forbidden",
	)
	ErrUnauthorized error = newPublic(
		"err_unauthorized",
	)
	ErrValidation error = errors.New(
		"validation",
	)
	ErrOTPCodeWrong error = newPublic(
		"err_otp_code_wrong",
	)
	ErrOTPSetupExpired error = newPublic(
		"err_otp_setup_expired",
	)
	ErrOTPTooManyAttempts error = newPublic(
		"err_otp_too_many_attempts",
	)
	ErrOAuthFailed error = newPublic(
		"err_oauth_failed",
	)
	ErrOAuthNoEmail error = newPublic(
		"err_oauth_no_email",
	)
	ErrOAuthEmailUnverified error = newPublic(
		"err_oauth_email_unverified",
	)
	ErrPollKindInvalid error = newPublic(
		"err_poll_kind_invalid",
	)
	ErrPollEnded error = newPublic(
		"err_poll_ended",
	)
	ErrPollSelectionInvalid error = newPublic(
		"err_poll_selection_invalid",
	)
	ErrPollNotAllowed error = newPublic(
		"err_poll_not_allowed",
	)
	ErrPollOptionsTooFew error = newPublic(
		"err_poll_options_too_few",
	)
	ErrPollOptionsTooMany error = newPublic(
		"err_poll_options_too_many",
	)
	ErrPollOptionTooLong error = newPublic(
		"err_poll_option_too_long",
	)
	ErrPollEndsAtInvalid error = newPublic(
		"err_poll_ends_at_invalid",
	)
	ErrPollEndsAtPast error = newPublic(
		"err_poll_ends_at_past",
	)
	ErrAuthenticationRequired error = newPublic(
		"err_authentication_required",
	)
	ErrOTPRequired error = newPublic(
		"err_otp_required",
	)
	ErrRateLimited error = newPublic(
		"err_rate_limited",
	)
	ErrInternal error = newPublic(
		"err_internal",
	)
	ErrAttachmentsDisabled error = newPublic(
		"err_attachments_disabled",
	)
	ErrReplyParentInvalid error = newPublic(
		"err_reply_parent_invalid",
	)
	ErrPicturesDisabled error = newPublic(
		"err_pictures_disabled",
	)
)
