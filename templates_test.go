package main

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/worker/targets/tmpl"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

var updateGolden = flag.Bool("update-golden", false,
	"rewrite testdata/golden from current template output")

func fixtureSystem() *setting.System {
	return &setting.System{
		Name:    "Example Board",
		BaseURL: "https://board.example.org",
	}
}

func fixtureSignup(t *testing.T) *signupconfirmation.SignupConfirmation {
	t.Helper()

	entity, err := signupconfirmation.New(
		&user.User{
			Username: "alice",
			Email:    "alice@example.org",
			Language: "en",
		},
		"Signup Confirmation",
		"IcCmCvjW",
		"session/signup",
	)
	if err != nil {
		t.Fatalf("signup fixture: %v", err)
	}

	entity.SetSystem(fixtureSystem())

	return entity
}

func fixtureReply(t *testing.T) *replynotification.ReplyNotification {
	t.Helper()

	entity, err := replynotification.New(
		&user.User{
			Username: "alice",
			Email:    "alice@example.org",
			Language: "en",
		},
		"New reply to a topic you follow",
	)
	if err != nil {
		t.Fatalf("reply fixture: %v", err)
	}

	rep := &reply.Reply{
		ID:      uuid.MustParse("019f64dc-7248-79f8-8ca7-b3cc23888087"),
		ShortID: "3Kf9vQ",
		Text:    "Have you tried turning it *off* and on again?",
		HTML:    "<p>Have you tried turning it <em>off</em> and on again?</p>",
	}

	entity.SetReply(rep, "bob", "_general/support/help#post-"+rep.ShortID)
	entity.SetCategory("General", "_general")
	entity.SetForum("Support", "_general/support")
	entity.SetTopic(
		uuid.MustParse("019f64dc-7248-79f8-8ca7-b3cc238880aa"),
		"Help",
		"_general/support/help",
	)

	entity.SetSystem(fixtureSystem())

	return entity
}

type goldenCase struct {
	spec    tmpl.Spec
	jobType asyncjob.JobType
	subType asyncjob.JobSubType
	lang    string
	data    func(*testing.T) any
}

func (c goldenCase) key() string {
	return strings.Join([]string{
		c.spec.Dir, string(c.jobType), string(c.subType), c.lang,
	}, "/")
}

func goldenCases() []goldenCase {
	signup := func(t *testing.T) any { return fixtureSignup(t) }
	reply := func(t *testing.T) any { return fixtureReply(t) }

	return []goldenCase{
		{tmpl.EmailSpec, asyncjob.Confirmation, asyncjob.Signup, "en", signup},
		{tmpl.EmailSpec, asyncjob.Notification, asyncjob.Reply, "en", reply},
		{tmpl.XMPPSpec, asyncjob.Confirmation, asyncjob.Signup, "en", signup},
		{tmpl.XMPPSpec, asyncjob.Notification, asyncjob.Reply, "en", reply},
	}
}

func goldenPath(c goldenCase, part string) string {
	return filepath.Join("testdata", "golden", c.spec.Dir,
		string(c.jobType)+"."+string(c.subType)+"."+c.lang+"."+part+".golden")
}

func assertGolden(t *testing.T, path string, got string) {
	t.Helper()

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (run: go test . -update-golden)", path, err)
	}

	if got != string(want) {
		t.Errorf("%s does not match rendered output\n--- got ---\n%s\n--- want ---\n%s",
			path, got, string(want))
	}
}

func TestTemplatesMatchGolden(t *testing.T) {
	for _, c := range goldenCases() {
		t.Run(c.key(), func(t *testing.T) {
			cache := tmpl.NewCache(&embedTemplates, c.spec)

			item, err := cache.TemplatesFor(c.jobType, c.subType, c.lang)
			if err != nil {
				t.Fatalf("load templates: %v", err)
			}

			data := c.data(t)

			var sb strings.Builder
			if err := item.TextTmpl.Execute(&sb, data); err != nil {
				t.Fatalf("render text: %v", err)
			}
			assertGolden(t, goldenPath(c, "text"), sb.String())

			if c.spec.HtmlExt == "" {
				if item.HtmlTmpl != nil {
					t.Error("HtmlTmpl is non-nil for a text-only spec")
				}
				return
			}

			sb.Reset()
			if err := item.HtmlTmpl.Execute(&sb, data); err != nil {
				t.Fatalf("render html: %v", err)
			}
			assertGolden(t, goldenPath(c, "html"), sb.String())
		})
	}
}

func TestEveryEmbeddedTemplateHasAGoldenCase(t *testing.T) {
	covered := make(map[string]bool)
	for _, c := range goldenCases() {
		covered[c.key()] = true
	}

	specs := map[string]tmpl.Spec{
		tmpl.EmailSpec.Dir: tmpl.EmailSpec,
		tmpl.XMPPSpec.Dir:  tmpl.XMPPSpec,
	}

	err := fs.WalkDir(&embedTemplates, "templates",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			rel := strings.TrimPrefix(path, "templates/")
			parts := strings.Split(rel, "/")
			if len(parts) != 3 {
				t.Errorf("unexpected template layout: %s", path)
				return nil
			}

			dir, jobType, file := parts[0], parts[1], parts[2]

			spec, ok := specs[dir]
			if !ok {
				t.Errorf("template %s lives under an unknown target dir %q", path, dir)
				return nil
			}

			ext := filepath.Ext(file)
			if ext != spec.TextExt && ext != spec.HtmlExt {
				t.Errorf("template %s has extension %q, which spec %q does not declare",
					path, ext, dir)
				return nil
			}

			name := strings.TrimSuffix(file, ".tmpl"+ext)
			subType, lang, found := strings.Cut(name, ".")
			if !found {
				t.Errorf("cannot parse subtype/lang out of %s", path)
				return nil
			}

			key := strings.Join([]string{dir, jobType, subType, lang}, "/")
			if !covered[key] {
				t.Errorf("template %s has no golden case (add %q to goldenCases)",
					path, key)
			}

			return nil
		})
	if err != nil {
		t.Fatalf("walk templates: %v", err)
	}
}
