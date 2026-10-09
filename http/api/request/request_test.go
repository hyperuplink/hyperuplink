package request

import (
	"reflect"
	"testing"
)

type bindInner struct {
	Text    string `json:"text"`
	TopicID string `json:"topic_id"`
}

type bindOuter struct {
	bindInner
	AttachmentIDs []string `json:"attachment_ids"`
	Untagged      string
	Skipped       string `json:"-"`
}

func TestJSONFieldNameFollowsTheNamespace(t *testing.T) {
	typ := reflect.TypeOf(&bindOuter{})

	for namespace, want := range map[string]string{
		"bindOuter.bindInner.TopicID": "topic_id",
		"bindOuter.bindInner.Text":    "text",
		"bindOuter.AttachmentIDs":     "attachment_ids",
		"bindOuter.Untagged":          "untagged",
		"bindOuter.Skipped":           "skipped",
		"bindOuter.Missing":           "missing",
	} {
		if got := jsonFieldName(typ, namespace); got != want {
			t.Errorf("jsonFieldName(%q) = %q, want %q", namespace, got, want)
		}
	}
}
