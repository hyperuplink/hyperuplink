package errs

import "errors"

var public = []error{
	ErrNoRows,
	ErrTargetIDNotFound,
	ErrUniqueViolationOn,
	ErrUnavailable,
}

func newPublic(text string) error {
	err := errors.New(text)
	public = append(public, err)

	return err
}

func Public(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrValidation) {
		return err
	}
	for _, known := range public {
		if errors.Is(err, known) {
			return known
		}
	}

	return ErrInternal
}
