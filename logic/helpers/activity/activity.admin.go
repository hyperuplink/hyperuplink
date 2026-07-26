package activity

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
)

func ChangedFields(before, after any) (changed []string) {
	bv := reflect.ValueOf(before)
	av := reflect.ValueOf(after)

	if bv.Kind() != reflect.Struct ||
		av.Kind() != reflect.Struct ||
		bv.Type() != av.Type() {
		return nil
	}

	t := bv.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if reflect.DeepEqual(
			bv.Field(i).Interface(),
			av.Field(i).Interface(),
		) {
			continue
		}

		changed = append(changed, fieldName(field))
	}

	return changed
}

func fieldName(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("json")
	if !ok {
		return field.Name
	}

	name := strings.Split(tag, ",")[0]
	if name == "" || name == "-" {
		return field.Name
	}

	return name
}

func RecordAdminVisit(
	rt *runtime.Runtime,
	actorID uuid.UUID,
	path string,
) {
	rec, err := activity.NewAdminVisit(actorID, path)
	if err != nil {
		rt.Error("error", err)
		return
	}

	if err = gh.Activity(rt).Record(rec); err != nil {
		rt.Error("error", err)
	}
}

func RecordAdminSettingsUpdate(
	rt *runtime.Runtime,
	actorID uuid.UUID,
	settingID string,
	before, after any,
) {
	changed := ChangedFields(before, after)
	if len(changed) == 0 {
		return
	}

	rec, err := activity.NewAdminSettingsUpdate(actorID, settingID, changed)
	if err != nil {
		rt.Error("error", err)
		return
	}

	if err = gh.Activity(rt).Record(rec); err != nil {
		rt.Error("error", err)
	}
}
