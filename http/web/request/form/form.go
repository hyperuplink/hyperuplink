package form

import (
	"reflect"

	"github.com/mrusme/hyperuplink/errs"
)

type Form struct {
	frm map[string]interface{}
}

func New() *Form {
	f := new(Form)

	return f
}

func (f *Form) Set(frm interface{}) error {
	f.frm = make(map[string]interface{})
	val := reflect.ValueOf(frm)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errs.ErrFormInvalid
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("form")
		if tag != "" {
			f.frm[tag] = fieldValue.Interface()
		} else {
			f.frm[field.Name] = fieldValue.Interface()
		}
	}
	return nil
}

func (f *Form) ValueFor(field string) interface{} {
	if val, exist := f.frm[field]; exist {
		return val
	}

	return nil
}

func (f *Form) StringValueFor(field string) string {
	val := f.ValueFor(field)
	if val != nil {
		return val.(string)
	}

	return ""
}
