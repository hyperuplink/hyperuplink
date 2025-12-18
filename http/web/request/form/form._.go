package form

import (
	"fmt"
	"html/template"
	"reflect"

	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/web/request/flash"
	"github.com/mrusme/hyperuplink/http/web/request/in"
)

type Form struct {
	fl  *flash.Flash
	in  *in.Internationalization
	frm map[string]interface{}
}

func New(fl *flash.Flash, in *in.Internationalization) *Form {
	f := new(Form)
	f.fl = fl
	f.in = in
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

func (f *Form) Input(
	inputType string,
	name string,
	value string,
	hasLabel bool,
	args ...string,
) template.HTML {
	var html string = ""

	if hasLabel {
		html = fmt.Sprintf("<label for=\"%s\">%s</label>",
			name, f.in.T(name))
	}

	fvalue := f.StringValueFor(name)
	if fvalue != "" {
		value = fvalue
	}

	var kv string = ""
	for i := 0; i < len(args); i++ {
		key := args[i]
		i++
		val := args[i]

		if val == "" {
			kv = fmt.Sprintf(`%s %s`, kv, key)
		} else {
			kv = fmt.Sprintf(`%s %s="%s"`, kv, key, val)
		}
	}

	html = fmt.Sprintf(
		`%s<input type="%s" name="%s" value="%s" class="%s" %s>`,
		html, inputType, name, value, f.fl.ClassFor(name), kv,
	)

	return template.HTML(html)
}
