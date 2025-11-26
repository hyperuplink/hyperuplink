package site

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (s *Site) ValidateForm(f any, t reflect.Type) (map[string]error, bool) {
	if err := s.c.Bind().Form(f); err != nil {
		var errs map[string]error = make(map[string]error)

		if valErrs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range valErrs {
				if t.Kind() != reflect.Struct {
					errs["error"] = valErrs
					break
				}

				field, ok := t.FieldByName(e.StructField())
				if !ok {
					errs["error"] = valErrs
					break
				}

				formTag, ok := field.Tag.Lookup("form")
				if !ok {
					errs["error"] = valErrs
					break
				}

				errs[formTag] = errors.New(s.T(fmt.Sprintf(
					"validation_%s_%s",
					strings.ToLower(e.Field()),
					e.Tag(),
				)))
			}
		} else {
			errs["error"] = err
		}

		return errs, false
	}

	return nil, true
}
