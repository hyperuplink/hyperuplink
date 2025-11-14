package errs

import "errors"

var (
	ErrConfigTypeUnsupported error = errors.New(
		"Config type unsupported",
	)

	ErrIfaceTypeUnsupported error = errors.New(
		"Interface type unsupported",
	)
)
