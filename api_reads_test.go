package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

var leakingKeys = map[string]bool{
	"author_email":      true,
	"author_member_of":  true,
	"email":             true,
	"email_unconfirmed": true,
	"member_of":         true,
	"otp_secret":        true,
	"password":          true,
	"secret_hash":       true,
}

func assertNoLeak(t *testing.T, body []byte) {
	t.Helper()

	var doc any
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatalf("not JSON: %v", err)
	}

	var walk func(path string, v any)
	walk = func(path string, v any) {
		switch x := v.(type) {
		case map[string]any:
			for k, child := range x {
				if leakingKeys[k] {
					t.Errorf("payload leaks %q at %s", k, path)
				}
				walk(path+"."+k, child)
			}
		case []any:
			for i, child := range x {
				walk(path+"[]", child)
				_ = i
			}
		}
	}
	walk("$", doc)
}

func items(t *testing.T, body map[string]any, key string) []map[string]any {
	t.Helper()

	raw, ok := body[key].([]any)
	if !ok {
		t.Fatalf("%q is not an array: %v", key, body[key])
	}

	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("%q holds a non-object: %v", key, item)
		}
		out = append(out, m)
	}

	return out
}

func object(t *testing.T, body map[string]any, key string) map[string]any {
	t.Helper()

	m, ok := body[key].(map[string]any)
	if !ok {
		t.Fatalf("%q is not an object: %v", key, body[key])
	}

	return m
}

func TestAPISessionDescribesTheViewer(t *testing.T) {
	h := apiTest(t)
	h.setDefaultPermission(t, permission.LevelRead)

	res := h.get(t, "/api/v1/session", "")
	expectStatus(t, res, 200)
	body := res.JSON(t)
	if body["authenticated"] != false || body["role"] != "guest" || body["user"] != nil {
		t.Errorf("guest session = %s", res.Body)
	}
	perms := object(t, body, "permissions")
	if perms["default"] != "read" || perms["admin"] != false {
		t.Errorf("guest permissions = %v", perms)
	}

	res = h.get(t, "/api/v1/session", h.keys["juno"])
	expectStatus(t, res, 200)
	body = res.JSON(t)
	if body["authenticated"] != true || body["role"] != "user" {
		t.Errorf("user session = %s", res.Body)
	}
	usr := object(t, body, "user")
	if usr["username"] != "juno" || usr["email"] != "juno@example.org" || usr["notify_on_reply"] != true {
		t.Errorf("user = %v", usr)
	}
	if groups := usr["groups"].([]any); len(groups) != 1 {
		t.Errorf("juno groups = %v", groups)
	}
	if key := object(t, body, "key"); key["name"] != "test" {
		t.Errorf("key = %v", key)
	}
	perms = object(t, body, "permissions")
	cats := object(t, perms, "categories")
	if perms["default"] != "read" || cats[h.categories["support"].String()] != "read_write" || cats[h.categories["general"].String()] != "read" {
		t.Errorf("juno permissions = %v", perms)
	}

	res = h.get(t, "/api/v1/session", h.keys["admin"])
	expectStatus(t, res, 200)
	perms = object(t, res.JSON(t), "permissions")
	if perms["admin"] != true || perms["default"] != "read_write_moderate" {
		t.Errorf("admin permissions = %v", perms)
	}
}

func TestAPIGuestsAreCappedAtRead(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/session", "")
	expectStatus(t, res, 200)
	if perms := object(t, res.JSON(t), "permissions"); perms["default"] != "read" {
		t.Errorf("a guest on a read-write board should resolve to read, got %v", perms["default"])
	}

	res = h.get(t, "/api/v1/forums", "")
	expectStatus(t, res, 200)
	for _, fum := range items(t, res.JSON(t), "forums") {
		p := object(t, fum, "permissions")
		if p["read"] != true || p["write"] != false {
			t.Errorf("guest permissions on %v = %v", fum["slug"], p)
		}
	}
}

func TestAPIBoardFollowsThePermissionDefault(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/board", "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)
	cats := items(t, body, "categories")
	if len(cats) != 2 || len(items(t, body, "recent_topics")) != 2 {
		t.Fatalf("board for a guest = %s", res.Body)
	}
	if forums := items(t, cats[0], "forums"); len(forums) != 1 || forums[0]["slug"] != "announcements" {
		t.Errorf("first category forums = %v", forums)
	}

	h.setDefaultPermission(t, permission.LevelNone)

	res = h.get(t, "/api/v1/board", "")
	expectStatus(t, res, 200)
	body = res.JSON(t)
	if len(items(t, body, "categories")) != 0 || len(items(t, body, "recent_topics")) != 0 {
		t.Errorf("board for a guest on a closed board = %s", res.Body)
	}

	res = h.get(t, "/api/v1/board", h.keys["juno"])
	expectStatus(t, res, 200)
	cats = items(t, res.JSON(t), "categories")
	if len(cats) != 1 || cats[0]["slug"] != "support" {
		t.Errorf("board for juno on a closed board = %s", res.Body)
	}
	if p := object(t, cats[0], "permissions"); p["write"] != true {
		t.Errorf("juno should write in support: %v", p)
	}

	res = h.get(t, "/api/v1/board", h.keys["admin"])
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "categories")) != 2 {
		t.Errorf("an administrator should see everything: %s", res.Body)
	}
}

func TestAPICategoriesByIDOrSlug(t *testing.T) {
	h := apiTest(t)

	bySlug := h.get(t, "/api/v1/categories/general", "")
	expectStatus(t, bySlug, 200)
	byID := h.get(t, "/api/v1/categories/"+h.categories["general"].String(), "")
	expectStatus(t, byID, 200)

	cat := object(t, bySlug.JSON(t), "category")
	if cat["id"] != h.categories["general"].String() || cat["id"] != object(t, byID.JSON(t), "category")["id"] {
		t.Errorf("category by slug and by id differ: %v", cat["id"])
	}
	if forums := items(t, cat, "forums"); len(forums) != 1 {
		t.Errorf("general forums = %v", forums)
	}
	if !strings.HasSuffix(cat["url"].(string), "/_general") {
		t.Errorf("category url = %v", cat["url"])
	}

	expectProblem(t, h.get(t, "/api/v1/categories/nope", ""), 404, "no_rows")

	res := h.get(t, "/api/v1/categories", "")
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "categories")) != 2 {
		t.Errorf("categories = %s", res.Body)
	}

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.get(t, "/api/v1/categories/general", ""), 404, "no_rows")
	expectStatus(t, h.get(t, "/api/v1/categories/support", h.keys["juno"]), 200)
	expectProblem(t, h.get(t, "/api/v1/categories/general", h.keys["juno"]), 404, "no_rows")
}

func TestAPIForumsListPermissionsPerViewer(t *testing.T) {
	h := apiTest(t)
	h.setDefaultPermission(t, permission.LevelRead)

	res := h.get(t, "/api/v1/forums", h.keys["juno"])
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)

	writable := map[string]bool{}
	for _, fum := range items(t, res.JSON(t), "forums") {
		writable[fum["slug"].(string)] = object(t, fum, "permissions")["write"] == true
		if cat := object(t, fum, "category"); cat["slug"] == "" || cat["id"] == "" {
			t.Errorf("forum %v has no category reference: %v", fum["slug"], cat)
		}
	}
	if writable["announcements"] || !writable["help"] {
		t.Errorf("juno writable forums = %v, want help only", writable)
	}

	res = h.get(t, "/api/v1/forums/help", h.keys["vera"])
	expectStatus(t, res, 200)
	fum := object(t, res.JSON(t), "forum")
	if fum["id"] != h.forums["help"].String() || object(t, fum, "permissions")["write"] != false {
		t.Errorf("forum help for vera = %v", fum)
	}

	expectStatus(t, h.get(t, "/api/v1/forums/"+h.forums["help"].String(), ""), 200)
	expectProblem(t, h.get(t, "/api/v1/forums/nope", ""), 404, "no_rows")
}

func TestAPIForumTopicsSortAndPaginate(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/forums/announcements/topics", "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)

	if fum := object(t, body, "forum"); fum["slug"] != "announcements" {
		t.Errorf("forum = %v", fum)
	}
	tops := items(t, body, "topics")
	if len(tops) != 1 {
		t.Fatalf("topics = %v", tops)
	}
	top := tops[0]
	if top["name"] != "Welcome to the board" || top["replies"] != float64(2) || top["kind"] != "regular" {
		t.Errorf("topic = %v", top)
	}
	if author := object(t, top, "author"); author["username"] != "admin" || author["role"] != "admin" {
		t.Errorf("author = %v", author)
	}
	if !strings.HasSuffix(top["url"].(string), "/_general/announcements/welcome-to-the-board") {
		t.Errorf("url = %v", top["url"])
	}
	for _, key := range []string{"created_at", "last_activity_at"} {
		if ts, _ := top[key].(string); !strings.HasSuffix(ts, "Z") {
			t.Errorf("%s = %v, want RFC 3339 in UTC", key, top[key])
		}
	}
	if top["viewer"] != nil {
		t.Errorf("a guest has no viewer state: %v", top["viewer"])
	}
	if pg := object(t, body, "pagination"); pg["total"] != float64(1) || pg["pages"] != float64(1) || pg["page"] != float64(1) {
		t.Errorf("pagination = %v", pg)
	}

	for _, sort := range []string{"active", "new", "replies"} {
		expectStatus(t, h.get(t, "/api/v1/forums/announcements/topics?sort="+sort, ""), 200)
	}
	body = expectProblem(t, h.get(t, "/api/v1/forums/announcements/topics?sort=bogus", ""), 422, "validation")
	if errs := object(t, body, "errors"); errs["sort"] == nil {
		t.Errorf("errors = %v", errs)
	}
	body = expectProblem(t, h.get(t, "/api/v1/forums/announcements/topics?page=0", ""), 422, "validation")
	if errs := object(t, body, "errors"); errs["page"] == nil {
		t.Errorf("errors = %v", errs)
	}
	body = expectProblem(t, h.get(t, "/api/v1/forums/announcements/topics?per_page=500&page=x", ""), 422, "validation")
	if errs := object(t, body, "errors"); errs["per_page"] == nil || errs["page"] == nil {
		t.Errorf("errors = %v", errs)
	}

	res = h.get(t, "/api/v1/forums/announcements/topics", h.keys["vera"])
	expectStatus(t, res, 200)
	if top := items(t, res.JSON(t), "topics")[0]; top["viewer"] == nil {
		t.Errorf("a signed-in viewer should get the read marker: %v", top)
	}
}

func TestAPITopicsFeedHonorsPermissions(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/topics?per_page=1", "")
	expectStatus(t, res, 200)
	body := res.JSON(t)
	if len(items(t, body, "topics")) != 1 || object(t, body, "pagination")["pages"] != float64(2) {
		t.Errorf("feed page one = %s", res.Body)
	}

	res = h.get(t, "/api/v1/topics?sort=new", "")
	expectStatus(t, res, 200)
	tops := items(t, res.JSON(t), "topics")
	if len(tops) != 2 || tops[0]["name"] != "Need help with my modem" {
		t.Errorf("newest first = %v", tops)
	}

	h.setDefaultPermission(t, permission.LevelNone)

	for token, want := range map[string]int{"": 0, h.keys["vera"]: 0, h.keys["juno"]: 1, h.keys["admin"]: 2} {
		res = h.get(t, "/api/v1/topics", token)
		expectStatus(t, res, 200)
		body = res.JSON(t)
		if got := len(items(t, body, "topics")); got != want || object(t, body, "pagination")["total"] != float64(want) {
			t.Errorf("token %q: %d topics, total %v, want %d", token, got, object(t, body, "pagination")["total"], want)
		}
	}
}

func TestAPITopicByIDOrShortID(t *testing.T) {
	h := apiTest(t)
	welcome := h.topics["welcome"]

	byID := h.get(t, "/api/v1/topics/"+welcome.ID.String(), "")
	expectStatus(t, byID, 200)
	assertNoLeak(t, byID.Body)
	byShort := h.get(t, "/api/v1/topics/"+welcome.ShortID, "")
	expectStatus(t, byShort, 200)

	top := object(t, byID.JSON(t), "topic")
	if top["id"] != welcome.ID.String() || top["id"] != object(t, byShort.JSON(t), "topic")["id"] {
		t.Errorf("topic by id and by short id differ")
	}
	if top["poll"] != nil || top["viewer"] != nil {
		t.Errorf("a regular topic for a guest has no poll and no viewer: %v", top)
	}
	if atts, _ := top["attachments"].([]any); atts == nil {
		t.Errorf("attachments should be an empty array, not null")
	}

	res := h.get(t, "/api/v1/topics/"+welcome.ShortID, h.keys["vera"])
	expectStatus(t, res, 200)
	if viewer := object(t, object(t, res.JSON(t), "topic"), "viewer"); viewer["unread"] == nil {
		t.Errorf("viewer = %v", viewer)
	}

	expectProblem(t, h.get(t, "/api/v1/topics/doesnotexist", ""), 404, "no_rows")

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.get(t, "/api/v1/topics/"+welcome.ShortID, ""), 404, "no_rows")
	expectProblem(t, h.get(t, "/api/v1/topics/"+welcome.ShortID, h.keys["juno"]), 404, "no_rows")
	expectStatus(t, h.get(t, "/api/v1/topics/"+h.topics["help"].ShortID, h.keys["juno"]), 200)
}

func TestAPIRepliesArePagedOldestFirst(t *testing.T) {
	h := apiTest(t)
	welcome := h.topics["welcome"]

	res := h.get(t, "/api/v1/topics/"+welcome.ShortID+"/replies", "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)

	reps := items(t, body, "replies")
	if len(reps) != 2 {
		t.Fatalf("replies = %v", reps)
	}
	if reps[0]["id"] != h.replies["first"].String() || reps[0]["parent_id"] != nil {
		t.Errorf("first reply = %v", reps[0])
	}
	if reps[1]["parent_id"] != h.replies["first"].String() {
		t.Errorf("second reply parent = %v", reps[1]["parent_id"])
	}
	if !strings.Contains(reps[1]["url"].(string), "/_general/announcements/welcome-to-the-board#post-") {
		t.Errorf("reply url = %v", reps[1]["url"])
	}
	if object(t, body, "pagination")["total"] != float64(2) {
		t.Errorf("pagination = %v", body["pagination"])
	}

	res = h.get(t, "/api/v1/topics/"+welcome.ShortID+"/replies?per_page=1&page=2", "")
	expectStatus(t, res, 200)
	body = res.JSON(t)
	if reps = items(t, body, "replies"); len(reps) != 1 || reps[0]["id"] != h.replies["second"].String() {
		t.Errorf("page two = %v", reps)
	}
	if pg := object(t, body, "pagination"); pg["pages"] != float64(2) || pg["page"] != float64(2) {
		t.Errorf("pagination = %v", pg)
	}
}

func TestAPIReplyShowsItsPosition(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/replies/"+h.replies["second"].String(), "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)

	if body["position"] != float64(2) || body["page"] != float64(1) {
		t.Errorf("position = %v, page = %v", body["position"], body["page"])
	}
	if top := object(t, body, "topic"); top["slug"] != "welcome-to-the-board" || top["id"] != h.topics["welcome"].ID.String() {
		t.Errorf("topic = %v", top)
	}
	rep := object(t, body, "reply")
	if rep["id"] != h.replies["second"].String() {
		t.Errorf("reply = %v", rep)
	}

	byShort := h.get(t, "/api/v1/replies/"+rep["short_id"].(string), "")
	expectStatus(t, byShort, 200)

	expectProblem(t, h.get(t, "/api/v1/replies/"+uuid.NewString(), ""), 404, "no_rows")

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.get(t, "/api/v1/replies/"+h.replies["second"].String(), ""), 404, "no_rows")
}

func TestAPISearchMatchesReadableContent(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/search?q=modem", "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)
	results := items(t, body, "results")
	if len(results) != 1 || results[0]["kind"] != "topic" || object(t, results[0], "topic")["slug"] != "need-help-with-my-modem" {
		t.Errorf("results = %v", results)
	}
	if object(t, body, "pagination")["total"] != float64(1) {
		t.Errorf("pagination = %v", body["pagination"])
	}

	expectStatus(t, h.get(t, "/api/v1/search?q=modem&in=title", ""), 200)
	res = h.get(t, "/api/v1/search?q=modem&in=replies", "")
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "results")) != 0 {
		t.Errorf("searching replies for a title word found %s", res.Body)
	}

	body = expectProblem(t, h.get(t, "/api/v1/search", ""), 422, "validation")
	if errs := object(t, body, "errors"); errs["q"] == nil {
		t.Errorf("errors = %v", errs)
	}
	expectProblem(t, h.get(t, "/api/v1/search?q=modem&in=bogus", ""), 422, "validation")

	h.setDefaultPermission(t, permission.LevelNone)
	res = h.get(t, "/api/v1/search?q=modem", "")
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "results")) != 0 {
		t.Errorf("a guest on a closed board found %s", res.Body)
	}
	res = h.get(t, "/api/v1/search?q=modem", h.keys["juno"])
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "results")) != 1 {
		t.Errorf("juno should find the support topic: %s", res.Body)
	}
}

func TestAPIUsersShowTheirOwnGroups(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/users/juno", "")
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	body := res.JSON(t)
	usr := object(t, body, "user")
	groups := items(t, usr, "groups")
	if len(groups) != 1 || groups[0]["id"] != "veterans" || groups[0]["name"] != "Veterans" {
		t.Errorf("juno groups = %v", groups)
	}
	if !strings.HasSuffix(usr["url"].(string), "/~juno") {
		t.Errorf("user url = %v", usr["url"])
	}

	res = h.get(t, "/api/v1/users/vera", "")
	expectStatus(t, res, 200)
	body = res.JSON(t)
	if groups := items(t, object(t, body, "user"), "groups"); len(groups) != 0 {
		t.Errorf("vera groups = %v", groups)
	}
	if tops := items(t, body, "topics"); len(tops) != 1 || tops[0]["slug"] != "need-help-with-my-modem" {
		t.Errorf("vera topics = %v", tops)
	}
	if reps := items(t, body, "replies"); len(reps) != 1 {
		t.Errorf("vera replies = %v", reps)
	}

	expectProblem(t, h.get(t, "/api/v1/users/nobody", ""), 404, "no_rows")
}

func TestAPIPagesFollowTheGeneralSettings(t *testing.T) {
	h := apiTest(t)

	expectProblem(t, h.get(t, "/api/v1/pages/about", ""), 404, "no_rows")
	expectProblem(t, h.get(t, "/api/v1/pages/nope", ""), 404, "no_rows")

	repo := gh.Repositories(h.rt).Setting
	general, err := settingRepo.GetByID[setting.General](h.ctx, repo, "general")
	if err != nil {
		t.Fatal(err)
	}
	general.JSONValue.EnableAbout = true
	general.JSONValue.About = "This board runs *Hyperuplink*."
	if err = settingRepo.Update(h.ctx, repo, general); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		general.JSONValue.EnableAbout = false
		if err := settingRepo.Update(h.ctx, repo, general); err != nil {
			t.Error(err)
		}
	})

	res := h.get(t, "/api/v1/pages/about", "")
	expectStatus(t, res, 200)
	page := object(t, res.JSON(t), "page")
	if page["name"] != "about" || !strings.Contains(page["html"].(string), "<em>Hyperuplink</em>") {
		t.Errorf("page = %v", page)
	}
}

func TestAPIUnknownAttachmentsAreNotFound(t *testing.T) {
	h := apiTest(t)

	expectProblem(t, h.get(t, "/api/v1/attachments/"+uuid.NewString(), ""), 404, "no_rows")
	expectProblem(t, h.get(t, "/api/v1/attachments/not-a-uuid", ""), 404, "no_rows")
}

func TestAPICategoryDefaultsOverrideTheBoardDefault(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]
	general := h.categories["general"].String()

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"level": "inherit",
	}), 422, "validation")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"group_id": "veterans", "category_id": h.categories["support"].String(), "level": "inherit",
	}), 422, "validation")

	h.setDefaultPermission(t, permission.LevelNone)
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"category_id": general, "level": "read",
	}), 204)
	t.Cleanup(func() {
		h.do(t, http.MethodPost, "/api/v1/admin/permissions/remove", admin, map[string]any{"category_id": general})
	})

	res := h.get(t, "/api/v1/board", "")
	expectStatus(t, res, 200)
	if cats := items(t, res.JSON(t), "categories"); len(cats) != 1 || cats[0]["slug"] != "general" {
		t.Errorf("guest board with one public category on a closed board = %s", res.Body)
	}
	expectProblem(t, h.get(t, "/api/v1/categories/support", ""), 404, "no_rows")
	res = h.get(t, "/api/v1/topics", "")
	expectStatus(t, res, 200)
	if object(t, res.JSON(t), "pagination")["total"] != float64(1) {
		t.Errorf("guest feed on a closed board with one public category = %s", res.Body)
	}
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": h.forums["announcements"].String(), "name": "Read only", "text": "x",
	}), 403, "err_forbidden")
	res = h.get(t, "/api/v1/board", h.keys["juno"])
	expectStatus(t, res, 200)
	if cats := items(t, res.JSON(t), "categories"); len(cats) != 2 {
		t.Errorf("juno should see general by its default and support by her group: %s", res.Body)
	}

	h.setDefaultPermission(t, permission.LevelReadWrite)
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"category_id": general, "level": "none",
	}), 204)

	res = h.get(t, "/api/v1/board", "")
	expectStatus(t, res, 200)
	if cats := items(t, res.JSON(t), "categories"); len(cats) != 1 || cats[0]["slug"] != "support" {
		t.Errorf("guest board with one closed category on an open board = %s", res.Body)
	}
	res = h.get(t, "/api/v1/topics", "")
	expectStatus(t, res, 200)
	body := res.JSON(t)
	for _, top := range items(t, body, "topics") {
		if object(t, top, "category")["slug"] == "general" {
			t.Errorf("a closed category's topic is in the feed: %v", top["name"])
		}
	}
	if object(t, body, "pagination")["total"] != float64(1) {
		t.Errorf("open board minus one category = %s", res.Body)
	}
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": h.forums["announcements"].String(), "name": "Hidden", "text": "x",
	}), 404, "no_rows")
	res = h.get(t, "/api/v1/forums/help", h.keys["vera"])
	expectStatus(t, res, 200)
	if p := object(t, object(t, res.JSON(t), "forum"), "permissions"); p["write"] != true {
		t.Errorf("vera should still write in the open category: %v", p)
	}

	res = h.get(t, "/api/v1/admin/permissions", admin)
	expectStatus(t, res, 200)
	levels := map[string]any{}
	for _, cat := range items(t, res.JSON(t), "categories") {
		levels[cat["slug"].(string)] = cat["level"]
	}
	if levels["general"] != "none" || levels["support"] != "inherit" {
		t.Errorf("category defaults = %v", levels)
	}

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"category_id": general, "level": "inherit",
	}), 204)
	res = h.get(t, "/api/v1/board", "")
	expectStatus(t, res, 200)
	if len(items(t, res.JSON(t), "categories")) != 2 {
		t.Errorf("inherit should remove the category default: %s", res.Body)
	}
}
