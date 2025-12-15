package errs

import "errors"

var (
	ErrNotImplemented error = errors.New(
		"Not implemented",
	)
	ErrConfigTypeUnsupported error = errors.New(
		"Config type unsupported",
	)
	ErrIfaceTypeUnsupported error = errors.New(
		"Interface type unsupported",
	)
	ErrTargetIDNotFound error = errors.New(
		"Target ID not found",
	)
	ErrStorageIDNotFound error = errors.New(
		"Storage ID not found",
	)
	ErrStorageTypeInvalid error = errors.New(
		"Storage type is invalid",
	)
	ErrFilePathInvalid error = errors.New(
		"File path is invalid",
	)
	ErrRedisAddrsEmpty error = errors.New(
		"Redis.Addresses cannot be empty",
	)
	ErrRedisAddrsMalformed error = errors.New(
		"Redis.Addresses malformed, needs to be <host>:<port>",
	)
	ErrHashInvalid error = errors.New(
		"Hash has incorrect format",
	)
	ErrHashVariantIncompatible error = errors.New(
		"Hash is incompatible variant",
	)
	ErrHashVersionIncompatible error = errors.New(
		"Hash is incompatible version",
	)
	ErrFormInvalid error = errors.New(
		"Form must be a struct or a pointer to a struct",
	)
	ErrJobTypeInvalid error = errors.New(
		"Job type is invalid",
	)
	ErrJobSubTypeInvalid error = errors.New(
		"Job sub type is invalid",
	)
	ErrJobPayloadInvalid error = errors.New(
		"Job payload is invalid",
	)
)
