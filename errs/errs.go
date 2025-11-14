package errs

import "errors"

var (
	ErrConfigTypeUnsupported error = errors.New(
		"Config type unsupported",
	)

	ErrIfaceTypeUnsupported error = errors.New(
		"Interface type unsupported",
	)

	ErrRedisAddrsEmpty error = errors.New(
		"Redis.Addresses cannot be empty",
	)
	ErrRedisAddrsMalformed error = errors.New(
		"Redis.Addresses malformed, needs to be <host>:<port>",
	)
)
