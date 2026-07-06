package errs

import "errors"

var (
	ErrNotImplemented error = errors.New(
		"err_not_implemented",
	)
	ErrUserIDNotFound error = errors.New(
		"err_user_id_not_found",
	)
	ErrPasswordWrong error = errors.New(
		"err_password_wrong",
	)
	ErrConfigTypeUnsupported error = errors.New(
		"err_config_type_unsupported",
	)
	ErrIfaceTypeUnsupported error = errors.New(
		"err_iface_type_unsupported",
	)
	ErrTargetIDNotFound error = errors.New(
		"err_target_id_not_found",
	)
	ErrStorageIDNotFound error = errors.New(
		"err_storage_id_not_found",
	)
	ErrStorageTypeInvalid error = errors.New(
		"err_storage_type_invalid",
	)
	ErrFilePathInvalid error = errors.New(
		"err_file_path_invalid",
	)
	ErrRedisAddrsEmpty error = errors.New(
		"err_redis_addrs_empty",
	)
	ErrRedisAddrsMalformed error = errors.New(
		"err_redis_addrs_malformed",
	)
	ErrHashInvalid error = errors.New(
		"err_hash_invalid",
	)
	ErrHashVariantIncompatible error = errors.New(
		"err_hash_variant_incompatible",
	)
	ErrHashVersionIncompatible error = errors.New(
		"err_hash_version_incompatible",
	)
	ErrFormInvalid error = errors.New(
		"err_form_invalid",
	)
	ErrJobTypeInvalid error = errors.New(
		"err_job_type_invalid",
	)
	ErrJobSubTypeInvalid error = errors.New(
		"err_job_sub_type_invalid",
	)
	ErrJobPayloadInvalid error = errors.New(
		"err_job_payload_invalid",
	)
	ErrNoRows error = errors.New(
		"no_rows",
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
	ErrUniqueViolationOn error = errors.New(
		"unique_violation_on",
	)
	ErrValidation error = errors.New(
		"validation",
	)
	ErrNoSuchTargetType error = errors.New(
		"err_no_such_target_type",
	)
	ErrOTPCodeWrong error = errors.New(
		"err_otp_code_wrong",
	)
	ErrOTPSetupExpired error = errors.New(
		"err_otp_setup_expired",
	)
)
