package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	logicapiforums "xn--gckvb8fzb.com/hyperuplink/logic/root/forums"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	logicreplies "xn--gckvb8fzb.com/hyperuplink/logic/root/replies"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

var (
	webOnce sync.Once
	webApp  *fiber.App
	webErr  error
)

func webTest(t *testing.T) (*apiHarness, *fiber.App) {
	t.Helper()

	h := apiTest(t)

	webOnce.Do(func() {
		webApp, webErr = h.newWebApp()
	})
	if webErr != nil {
		t.Fatalf("web harness: %v", webErr)
	}

	return h, webApp
}

func (h *apiHarness) newWebApp() (*fiber.App, error) {
	cfg, err := gh.WebConfigFrom(h.rt.Config())
	if err != nil {
		return nil, err
	}

	srv, err := web.New(h.rt, cfg)
	if err != nil {
		return nil, err
	}
	if err = srv.Startup(); err != nil {
		return nil, err
	}

	return srv.App(), nil
}

func (h *apiHarness) webGet(t *testing.T, app *fiber.App, path string) (int, string) {
	t.Helper()

	res, err := app.Test(httptest.NewRequest(http.MethodGet, path, nil),
		fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("GET %s: read body: %v", path, err)
	}

	return res.StatusCode, string(body)
}

func TestWebForumAndTopicPathsMatchEverySegment(t *testing.T) {
	h, app := webTest(t)

	welcome := h.topics["welcome"]
	help := h.topics["help"]

	for _, tc := range []struct {
		path   string
		status int
		shows  string
	}{
		{"/_general/announcements", 200, "Announcements"},
		{"/_support/help", 200, "Help"},
		{"/_support/announcements", 404, ""},
		{"/_general/help", 404, ""},
		{"/_nope/announcements", 404, ""},
		{"/_general/announcements/" + welcome.Slug, 200, welcome.Name},
		{"/_support/help/" + help.Slug, 200, help.Name},
		{"/_support/announcements/" + welcome.Slug, 404, ""},
		{"/_general/help/" + help.Slug, 404, ""},
		{"/_support/help/" + welcome.Slug, 404, ""},
	} {
		status, body := h.webGet(t, app, tc.path)
		if status != tc.status {
			t.Errorf("GET %s = %d, want %d", tc.path, status, tc.status)
		}
		if tc.shows != "" && !strings.Contains(body, tc.shows) {
			t.Errorf("GET %s does not show %q", tc.path, tc.shows)
		}
	}
}

func TestWebShowFunctionsNeedTheMatchingCategory(t *testing.T) {
	h := apiTest(t)

	admin := h.users["admin"]
	perms := h.perms(admin)
	welcome := h.topics["welcome"]

	_, err := logicforums.Show(h.ctx, h.rt, &logicforums.ShowInput{
		CategorySlug: "support",
		ForumSlug:    "announcements",
		Page:         1,
		PerPage:      10,
	}, perms)
	if !errors.Is(err, errs.ErrNoRows) {
		t.Errorf("forum under the wrong category: err = %v, want ErrNoRows", err)
	}

	fview, err := logicforums.Show(h.ctx, h.rt, &logicforums.ShowInput{
		CategorySlug: "general",
		ForumSlug:    "announcements",
		Page:         1,
		PerPage:      10,
	}, perms)
	if err != nil {
		t.Fatalf("forum under its category: %v", err)
	}
	if fview.Forum.Slug != "announcements" || fview.Forum.CategorySlug != "general" {
		t.Errorf("forum under its category = %s in %s", fview.Forum.Slug, fview.Forum.CategorySlug)
	}

	_, err = logictopics.Show(h.ctx, h.rt, &logictopics.ShowInput{
		CategorySlug: "support",
		ForumSlug:    "announcements",
		TopicSlug:    welcome.Slug,
		Page:         1,
		PerPage:      10,
	}, perms)
	if !errors.Is(err, errs.ErrNoRows) {
		t.Errorf("topic under the wrong category: err = %v, want ErrNoRows", err)
	}

	tview, err := logictopics.Show(h.ctx, h.rt, &logictopics.ShowInput{
		CategorySlug: "general",
		ForumSlug:    "announcements",
		TopicSlug:    welcome.Slug,
		Page:         1,
		PerPage:      10,
	}, perms)
	if err != nil {
		t.Fatalf("topic under its category: %v", err)
	}
	if tview.Topic.ID != welcome.ID {
		t.Errorf("topic under its category = %s, want %s", tview.Topic.ID, welcome.ID)
	}

	err = logictopics.PollVoteBySlugs(h.ctx, h.rt, perms, &logictopics.PollVoteBySlugsInput{
		CategorySlug: "support",
		ForumSlug:    "announcements",
		TopicSlug:    welcome.Slug,
		AuthorID:     admin.ID,
		Selection:    1,
	})
	if !errors.Is(err, errs.ErrNoRows) {
		t.Errorf("poll vote under the wrong category: err = %v, want ErrNoRows", err)
	}
}

func TestWebRepliesWithOneTimestampKeepTheirOrderAcrossPages(t *testing.T) {
	h := apiTest(t)

	admin := h.users["admin"]
	perms := h.perms(admin)

	top, err := logicnewpost.Create(h.ctx, h.rt, admin.ID, perms, time.UTC, &logicnewpost.CreateInput{
		Name:    "Seven replies in one second",
		Text:    "Every reply below shares a timestamp.",
		ForumID: h.forums["announcements"].String(),
		Kind:    "regular",
	}, nil)
	if err != nil {
		t.Fatalf("topic: %v", err)
	}

	const replies = 7
	want := make([]uuid.UUID, 0, replies)
	for i := range replies {
		created, cerr := logictopics.CreateReply(h.ctx, h.rt, admin.ID, perms, &logictopics.CreateReplyInput{
			Text:    fmt.Sprintf("Reply %d", i+1),
			TopicID: top.ID.String(),
		}, nil)
		if cerr != nil {
			t.Fatalf("reply %d: %v", i+1, cerr)
		}
		want = append(want, created.Reply.ID)
	}
	slices.SortFunc(want, func(a, b uuid.UUID) int { return bytes.Compare(a[:], b[:]) })

	stamp := time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC)
	if _, err = database.From(h.rt).Exec(h.ctx,
		`UPDATE replies SET created_at = $1, updated_at = $1 WHERE topic_id = $2`,
		stamp, top.ID,
	); err != nil {
		t.Fatalf("one timestamp: %v", err)
	}

	got := make([]uuid.UUID, 0, replies)
	for page := 1; ; page++ {
		view, verr := logictopics.Show(h.ctx, h.rt, &logictopics.ShowInput{
			CategorySlug: "general",
			ForumSlug:    "announcements",
			TopicSlug:    top.Slug,
			Page:         page,
			PerPage:      3,
		}, perms)
		if verr != nil {
			t.Fatalf("page %d: %v", page, verr)
		}
		for _, rep := range view.Replies {
			got = append(got, rep.ID)
		}
		if page >= view.Pages {
			break
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("web pages = %v\nwant %v", got, want)
	}

	api, err := logicreplies.Index(h.ctx, h.rt, &logicreplies.IndexInput{Topic: top, Page: 1, PerPage: replies})
	if err != nil {
		t.Fatalf("api page: %v", err)
	}
	for i, rep := range api.Replies {
		if rep.ID != want[i] {
			t.Errorf("api reply %d = %s, want %s", i+1, rep.ID, want[i])
		}
	}

	for i, id := range want {
		position, perr := gh.Repositories(h.rt).Reply.PositionForUUID(h.ctx, id)
		if perr != nil {
			t.Fatalf("position of %s: %v", id, perr)
		}
		if position != int64(i+1) {
			t.Errorf("position of reply %d = %d", i+1, position)
		}
	}
}

func TestWebPagesAreNotServedFromAResponseCache(t *testing.T) {
	_, app := webTest(t)

	var cookies []string
	for i := range 3 {
		req := httptest.NewRequest(http.MethodGet, "/_general/announcements", nil)
		for _, cookie := range cookies {
			req.Header.Add(fiber.HeaderCookie, cookie)
		}
		res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
		if err != nil {
			t.Fatalf("request %d: %v", i+1, err)
		}
		res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatalf("request %d = %d", i+1, res.StatusCode)
		}
		for _, name := range []string{"X-Cache", fiber.HeaderAge} {
			if value := res.Header.Get(name); value != "" {
				t.Errorf("request %d: %s = %q, the page went through a response cache", i+1, name, value)
			}
		}
		if cc := res.Header.Get(fiber.HeaderCacheControl); strings.Contains(cc, "public") {
			t.Errorf("request %d: Cache-Control = %q", i+1, cc)
		}
		if i == 0 {
			for _, set := range res.Header.Values(fiber.HeaderSetCookie) {
				cookies = append(cookies, strings.SplitN(set, ";", 2)[0])
			}
			if len(cookies) == 0 {
				t.Fatal("the first response set no cookie")
			}
		}
	}
}

func TestWebForumPagesOrderTopicsLikeTheAPI(t *testing.T) {
	h := apiTest(t)

	admin := h.users["admin"]
	perms := h.perms(admin)

	const topics = 5
	created := make([]uuid.UUID, 0, topics)
	for i := range topics {
		top, err := logicnewpost.Create(h.ctx, h.rt, admin.ID, perms, time.UTC, &logicnewpost.CreateInput{
			Name:    fmt.Sprintf("Same moment %d", i+1),
			Text:    "Five topics, one timestamp, one of them pinned.",
			ForumID: h.forums["announcements"].String(),
			Kind:    "regular",
		}, nil)
		if err != nil {
			t.Fatalf("topic %d: %v", i+1, err)
		}
		created = append(created, top.ID)
	}
	t.Cleanup(func() {
		if _, err := database.From(h.rt).Exec(h.ctx,
			`UPDATE topics SET deleted_at = NOW(), pinned = FALSE WHERE id = ANY($1)`, created,
		); err != nil {
			t.Errorf("remove the topics: %v", err)
		}
	})

	stamp := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	if _, err := database.From(h.rt).Exec(h.ctx,
		`UPDATE topics SET created_at = $1, updated_at = $1 WHERE id = ANY($2)`, stamp, created,
	); err != nil {
		t.Fatal(err)
	}
	pinned := created[2]
	if _, err := database.From(h.rt).Exec(h.ctx,
		`UPDATE topics SET pinned = TRUE WHERE id = $1`, pinned,
	); err != nil {
		t.Fatal(err)
	}

	forum, err := logicapiforums.Show(h.ctx, h.rt, "announcements", perms)
	if err != nil {
		t.Fatal(err)
	}
	api, err := logicapiforums.ListTopics(h.ctx, h.rt, &logicapiforums.TopicsInput{
		Forum:   forum,
		Sort:    topic.SortActive,
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := make([]uuid.UUID, 0, len(api.Topics))
	for _, top := range api.Topics {
		want = append(want, top.ID)
	}
	if len(want) < topics || want[0] != pinned {
		t.Fatalf("the API lists %d topics, first %s, want the pinned %s first", len(want), want[0], pinned)
	}

	got := make([]uuid.UUID, 0, len(want))
	for page := 1; ; page++ {
		view, verr := logicforums.Show(h.ctx, h.rt, &logicforums.ShowInput{
			CategorySlug: "general",
			ForumSlug:    "announcements",
			Page:         page,
			PerPage:      2,
		}, perms)
		if verr != nil {
			t.Fatalf("page %d: %v", page, verr)
		}
		for _, top := range view.Topics {
			got = append(got, top.ID)
		}
		if page >= view.Pages {
			break
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("web pages = %v\napi list  = %v", got, want)
	}

	ties := make([]uuid.UUID, 0, topics-1)
	for _, id := range got {
		if id != pinned && slices.Contains(created, id) {
			ties = append(ties, id)
		}
	}
	sorted := slices.Clone(ties)
	slices.SortFunc(sorted, func(a, b uuid.UUID) int { return bytes.Compare(b[:], a[:]) })
	if !slices.Equal(ties, sorted) {
		t.Errorf("topics with one timestamp = %v, want id descending %v", ties, sorted)
	}
}

func TestWebShowsNoRawErrorText(t *testing.T) {
	h, app := webTest(t)

	if _, err := database.From(h.rt).Exec(h.ctx,
		`ALTER TABLE users ADD CONSTRAINT f38_probe CHECK (username <> 'rawerror')`,
	); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := database.From(h.rt).Exec(h.ctx,
			`ALTER TABLE users DROP CONSTRAINT f38_probe`,
		); err != nil {
			t.Errorf("drop the probe constraint: %v", err)
		}
	})

	wc := newWebClient(t, app, "en-US,en;q=0.9")
	_, page := wc.get("/session/signup")
	token := csrfField.FindStringSubmatch(page)
	if token == nil {
		t.Fatal("the signup page has no CSRF token")
	}

	res := wc.do(http.MethodPost, "/session/signup", url.Values{
		"username":        {"rawerror"},
		"email":           {"rawerror@example.org"},
		"password":        {testPassword},
		"password_repeat": {testPassword},
		"_csrf":           {token[1]},
	})
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)

	for _, leak := range []string{"SQLSTATE", "f38_probe", "violates check constraint"} {
		if strings.Contains(body, leak) {
			t.Errorf("the signup page shows %q", leak)
		}
	}
	if !strings.Contains(body, "Something went wrong. Please try again.") {
		t.Error("the signup page doesn't show the err_internal message")
	}
	if _, err = gh.Repositories(h.rt).User.GetByUsername(h.ctx, "rawerror", queryOne()); !errors.Is(err, errs.ErrNoRows) {
		t.Errorf("the failed signup created a user: %v", err)
	}
}
