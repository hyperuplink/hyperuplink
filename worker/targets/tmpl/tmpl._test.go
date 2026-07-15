package tmpl

import (
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"templates/email/confirmation/signup.en.tmpl.eml": {
			Data: []byte("email confirmation signup en text"),
		},
		"templates/email/confirmation/signup.en.tmpl.html": {
			Data: []byte("email confirmation signup en html"),
		},
		"templates/email/notification/reply.en.tmpl.eml": {
			Data: []byte("email notification reply en text"),
		},
		"templates/email/notification/reply.en.tmpl.html": {
			Data: []byte("email notification reply en html"),
		},
		"templates/email/confirmation/signup.de.tmpl.eml": {
			Data: []byte("email confirmation signup de text"),
		},
		"templates/email/confirmation/signup.de.tmpl.html": {
			Data: []byte("email confirmation signup de html"),
		},
		"templates/xmpp/confirmation/signup.en.tmpl.md": {
			Data: []byte("xmpp confirmation signup en markdown"),
		},
		"templates/xmpp/notification/reply.en.tmpl.md": {
			Data: []byte("xmpp notification reply en markdown"),
		},
	}
}

func renderText(t *testing.T, item Item) string {
	t.Helper()

	var sb strings.Builder
	if err := item.TextTmpl.Execute(&sb, nil); err != nil {
		t.Fatalf("execute text template: %v", err)
	}

	return sb.String()
}

func renderHtml(t *testing.T, item Item) string {
	t.Helper()

	var sb strings.Builder
	if err := item.HtmlTmpl.Execute(&sb, nil); err != nil {
		t.Fatalf("execute html template: %v", err)
	}

	return sb.String()
}

func TestSubTypesDoNotShareCacheEntry(t *testing.T) {
	c := NewCache(testFS(), EmailSpec)

	signup, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "en")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}

	reply, err := c.TemplatesFor(asyncjob.Notification, asyncjob.Reply, "en")
	if err != nil {
		t.Fatalf("reply: %v", err)
	}

	if got, want := renderText(t, signup), "email confirmation signup en text"; got != want {
		t.Errorf("signup text = %q, want %q", got, want)
	}

	if got, want := renderText(t, reply), "email notification reply en text"; got != want {
		t.Errorf("reply text = %q, want %q", got, want)
	}

	if got, want := renderHtml(t, reply), "email notification reply en html"; got != want {
		t.Errorf("reply html = %q, want %q", got, want)
	}
}

func TestLangsDoNotShareCacheEntry(t *testing.T) {
	c := NewCache(testFS(), EmailSpec)

	en, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "en")
	if err != nil {
		t.Fatalf("en: %v", err)
	}

	de, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "de")
	if err != nil {
		t.Fatalf("de: %v", err)
	}

	if got, want := renderText(t, en), "email confirmation signup en text"; got != want {
		t.Errorf("en text = %q, want %q", got, want)
	}

	if got, want := renderText(t, de), "email confirmation signup de text"; got != want {
		t.Errorf("de text = %q, want %q", got, want)
	}
}

func TestTextOnlySpecLeavesHtmlNil(t *testing.T) {
	c := NewCache(testFS(), XMPPSpec)

	item, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "en")
	if err != nil {
		t.Fatalf("signup: %v", err)
	}

	if item.HtmlTmpl != nil {
		t.Error("HtmlTmpl = non-nil, want nil for a text-only spec")
	}

	if got, want := renderText(t, item), "xmpp confirmation signup en markdown"; got != want {
		t.Errorf("signup text = %q, want %q", got, want)
	}
}

func TestSecondCallIsCached(t *testing.T) {
	fsys := testFS()
	c := NewCache(fsys, EmailSpec)

	first, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "en")
	if err != nil {
		t.Fatalf("first: %v", err)
	}

	delete(fsys, "templates/email/confirmation/signup.en.tmpl.eml")
	delete(fsys, "templates/email/confirmation/signup.en.tmpl.html")

	second, err := c.TemplatesFor(asyncjob.Confirmation, asyncjob.Signup, "en")
	if err != nil {
		t.Fatalf("second: %v", err)
	}

	if first.TextTmpl != second.TextTmpl {
		t.Error("second call reparsed the text template, want cached instance")
	}
}

func TestMissingTemplateErrors(t *testing.T) {
	c := NewCache(testFS(), EmailSpec)

	if _, err := c.TemplatesFor(
		asyncjob.Confirmation, asyncjob.EmailChange, "en",
	); err == nil {
		t.Error("err = nil, want an error for a missing template")
	}
}

func TestConcurrentTemplatesFor(t *testing.T) {
	c := NewCache(testFS(), EmailSpec)

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := c.TemplatesFor(
				asyncjob.Confirmation, asyncjob.Signup, "en",
			); err != nil {
				t.Errorf("TemplatesFor: %v", err)
			}
		}()
	}
	wg.Wait()
}
