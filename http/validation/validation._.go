package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	VALID_SLUG_REGEX     = `^[a-z0-9\-]+$`
	VALID_SLUG_REGEX_NEG = `[^a-z0-9\-]`
)

func IsValidSlug(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	re := regexp.MustCompile(VALID_SLUG_REGEX)

	return re.MatchString(value)
}

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("slug", IsValidSlug)

	return validate
}

type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator(validate *validator.Validate) *StructValidator {
	return &StructValidator{validate: validate}
}

func (v *StructValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
