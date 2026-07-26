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
	ErrUserIDNotFound error = errors.New(
		"err_user_id_not_found",
	)
	ErrPasswordWrong error = errors.New(
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
	ErrFormInvalid error = errors.New(
		"err_form_invalid",
	)
	ErrUsernamePasswordWrong error = errors.New(
		"username_password_wrong",
	)
	ErrEmailConfirmationTokenWrong error = errors.New(
		"email_confirmation_token_wrong",
	)
	ErrPictureTooLarge error = errors.New(
		"picture_too_large",
	)
	ErrPictureFormatNotAllowed error = errors.New(
		"picture_format_not_allowed",
	)
	ErrInvalidUploadFormat error = errors.New(
		"invalid_upload_format",
	)
	ErrInvalidStorageProvider error = errors.New(
		"invalid_storage_provider",
	)
	ErrThemeStorageNotConfigured error = errors.New(
		"theme_storage_not_configured",
	)
	ErrInvalidTheme error = errors.New(
		"invalid_theme",
	)
	ErrInvalidColorscheme error = errors.New(
		"invalid_colorscheme",
	)
	ErrBannerUploadFailed error = errors.New(
		"banner_upload_failed",
	)
	ErrBannerAlreadySet error = errors.New(
		"banner_already_set",
	)
	ErrBannerFormatNotAllowed error = errors.New(
		"banner_format_not_allowed",
	)
	ErrFaviconUploadFailed error = errors.New(
		"favicon_upload_failed",
	)
	ErrFaviconAlreadySet error = errors.New(
		"favicon_already_set",
	)
	ErrFaviconFormatNotAllowed error = errors.New(
		"favicon_format_not_allowed",
	)
	ErrBackgroundUploadFailed error = errors.New(
		"background_upload_failed",
	)
	ErrBackgroundAlreadySet error = errors.New(
		"background_already_set",
	)
	ErrBackgroundFormatNotAllowed error = errors.New(
		"background_format_not_allowed",
	)
	ErrAttachmentTooLarge error = errors.New(
		"attachment_too_large",
	)
	ErrAttachmentUploadFailed error = errors.New(
		"attachment_upload_failed",
	)
	ErrAttachmentFormatNotAllowed error = errors.New(
		"attachment_format_not_allowed",
	)
	ErrAttachmentHookFailed error = errors.New(
		"attachment_hook_failed",
	)
	ErrAttachmentDuplicate error = errors.New(
		"attachment_duplicate",
	)
	ErrAttachmentNotFound error = errors.New(
		"attachment_not_found",
	)
	ErrAPIKeyInvalid error = errors.New(
		"err_apikey_invalid",
	)
	ErrForbidden error = errors.New(
		"err_forbidden",
	)
	ErrUnauthorized error = errors.New(
		"err_unauthorized",
	)
	ErrValidation error = errors.New(
		"validation",
	)
	ErrOTPCodeWrong error = errors.New(
		"err_otp_code_wrong",
	)
	ErrOTPSetupExpired error = errors.New(
		"err_otp_setup_expired",
	)
	ErrOTPTooManyAttempts error = errors.New(
		"err_otp_too_many_attempts",
	)
	ErrOAuthFailed error = errors.New(
		"err_oauth_failed",
	)
	ErrOAuthNoEmail error = errors.New(
		"err_oauth_no_email",
	)
	ErrOAuthEmailUnverified error = errors.New(
		"err_oauth_email_unverified",
	)
	ErrPollKindInvalid error = errors.New(
		"err_poll_kind_invalid",
	)
	ErrPollEnded error = errors.New(
		"err_poll_ended",
	)
	ErrPollSelectionInvalid error = errors.New(
		"err_poll_selection_invalid",
	)
	ErrPollNotAllowed error = errors.New(
		"err_poll_not_allowed",
	)
	ErrPollOptionsTooFew error = errors.New(
		"err_poll_options_too_few",
	)
	ErrPollOptionsTooMany error = errors.New(
		"err_poll_options_too_many",
	)
	ErrPollOptionTooLong error = errors.New(
		"err_poll_option_too_long",
	)
	ErrPollEndsAtInvalid error = errors.New(
		"err_poll_ends_at_invalid",
	)
	ErrPollEndsAtPast error = errors.New(
		"err_poll_ends_at_past",
	)
)
