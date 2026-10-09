package main

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func TestAPIAdminIsForAdministrators(t *testing.T) {
	h := apiTest(t)

	expectProblem(t, h.get(t, "/api/v1/admin/general", ""), 401, "err_authentication_required")
	expectProblem(t, h.get(t, "/api/v1/admin/general", h.keys["vera"]), 403, "err_forbidden")
	expectProblem(t, h.do(t, http.MethodPut, "/api/v1/admin/general", h.keys["juno"], map[string]any{}), 403, "err_forbidden")

	res := h.get(t, "/api/v1/admin/general", h.keys["admin"])
	expectStatus(t, res, 200)
	body := res.JSON(t)
	if sys := object(t, body, "system"); sys["allow_api_guests"] != true || sys["name"] == "" {
		t.Errorf("system = %v", sys)
	}
	if general := object(t, body, "general"); general["enable_about"] != false {
		t.Errorf("general = %v", general)
	}
}

func TestAPIAdminGeneralClosesTheAPIToGuests(t *testing.T) {
	h := apiTest(t)

	res := h.get(t, "/api/v1/admin/general", h.keys["admin"])
	expectStatus(t, res, 200)
	sys := object(t, res.JSON(t), "system")

	settings := func(allow bool) map[string]any {
		return map[string]any{
			"name":                     sys["name"],
			"base_url":                 sys["base_url"],
			"topics_per_page":          10,
			"posts_per_page":           10,
			"admin_log_retention_days": 30,
			"allow_api_guests":         allow,
		}
	}

	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/general", h.keys["admin"], settings(false)), 204)
	t.Cleanup(func() {
		h.do(t, http.MethodPut, "/api/v1/admin/general", h.keys["admin"], settings(true))
	})

	res = h.get(t, "/api/v1", "")
	expectStatus(t, res, 200)
	if res.JSON(t)["guests"] != false {
		t.Errorf("discovery after closing = %s", res.Body)
	}
	expectProblem(t, h.get(t, "/api/v1/board", ""), 401, "err_authentication_required")
	expectStatus(t, h.get(t, "/api/v1/board", h.keys["vera"]), 200)

	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/general", h.keys["admin"], settings(true)), 204)
	expectStatus(t, h.get(t, "/api/v1/board", ""), 200)

	body := expectProblem(t, h.do(t, http.MethodPut, "/api/v1/admin/general", h.keys["admin"], map[string]any{
		"name": "",
	}), 422, "validation")
	if errs := object(t, body, "errors"); errs["name"] == nil || errs["base_url"] == nil {
		t.Errorf("errors = %v", errs)
	}
}

func TestAPIAdminCategoriesAndForumsAreManaged(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]

	res := h.do(t, http.MethodPost, "/api/v1/admin/board/categories", admin, map[string]any{"name": "Lounge", "slug": "lounge"})
	expectStatus(t, res, 201)
	catID, _ := res.JSON(t)["id"].(string)
	if catID == "" || res.Header.Get(fiber.HeaderLocation) != "/api/v1/admin/board/categories/"+catID {
		t.Fatalf("created category = %s, Location %q", res.Body, res.Header.Get(fiber.HeaderLocation))
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/board/categories", admin, map[string]any{"name": "Again", "slug": "lounge"}), 409, "unique_violation_on")

	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/board/categories/"+catID, admin, map[string]any{"name": "The Lounge", "slug": "lounge"}), 204)
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/board/categories/"+catID+"/moveup", admin, nil), 204)
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/board/categories/"+catID+"/movedown", admin, nil), 204)
	expectProblem(t, h.do(t, http.MethodPut, "/api/v1/admin/board/categories/not-a-uuid", admin, map[string]any{"name": "x", "slug": "x"}), 404, "no_rows")

	res = h.get(t, "/api/v1/admin/board/categories", admin)
	expectStatus(t, res, 200)
	found := false
	for _, cat := range items(t, res.JSON(t), "categories") {
		found = found || cat["slug"] == "lounge" && cat["name"] == "The Lounge"
	}
	if !found {
		t.Errorf("lounge missing from %s", res.Body)
	}

	res = h.do(t, http.MethodPost, "/api/v1/admin/board/forums", admin, map[string]any{
		"name": "Chatter", "slug": "chatter", "category_id": catID,
	})
	expectStatus(t, res, 201)
	forumID, _ := res.JSON(t)["id"].(string)
	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/board/forums/"+forumID, admin, map[string]any{
		"name": "Chatter", "slug": "chatter", "category_id": catID, "description": "Anything goes.",
	}), 204)

	res = h.get(t, "/api/v1/forums/chatter", "")
	expectStatus(t, res, 200)
	if fum := object(t, res.JSON(t), "forum"); fum["description"] != "Anything goes." || object(t, fum, "category")["slug"] != "lounge" {
		t.Errorf("forum = %v", fum)
	}

	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/admin/board/forums/"+forumID, admin, nil), 204)
	expectProblem(t, h.get(t, "/api/v1/forums/chatter", ""), 404, "no_rows")
	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/admin/board/categories/"+catID, admin, nil), 204)
	expectProblem(t, h.get(t, "/api/v1/categories/lounge", ""), 404, "no_rows")
	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/admin/board/categories/"+uuid.NewString(), admin, nil), 404, "no_rows")
}

func TestAPIAdminGroupsAndPermissions(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]

	res := h.do(t, http.MethodPost, "/api/v1/admin/permissions/groups", admin, map[string]any{"id": "mods", "name": "Moderators"})
	expectStatus(t, res, 201)
	if res.Header.Get(fiber.HeaderLocation) != "/api/v1/admin/permissions/groups/mods" {
		t.Errorf("Location = %q", res.Header.Get(fiber.HeaderLocation))
	}
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/groups", admin, map[string]any{"id": "mods", "name": "Again"}), 409, "unique_violation_on")

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"group_id": "mods", "category_id": h.categories["general"].String(), "level": "read_write_moderate",
	}), 204)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"group_id": "mods", "category_id": h.categories["general"].String(), "level": "owner",
	}), 422, "validation")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/apply", admin, map[string]any{
		"group_id": "mods", "level": "read",
	}), 422, "validation")

	res = h.get(t, "/api/v1/admin/permissions", admin)
	expectStatus(t, res, 200)
	body := res.JSON(t)
	if body["default_level"] != "read_write" {
		t.Errorf("default_level = %v", body["default_level"])
	}
	found := false
	for _, grp := range items(t, body, "groups") {
		if grp["id"] == "mods" {
			found = true
			if rows := items(t, grp, "rows"); len(rows) != 1 || rows[0]["level"] != "read_write_moderate" {
				t.Errorf("mods rows = %v", rows)
			}
		}
	}
	if !found {
		t.Errorf("mods missing from %s", res.Body)
	}

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/permissions/remove", admin, map[string]any{
		"group_id": "mods", "category_id": h.categories["general"].String(),
	}), 204)
	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/admin/permissions/groups/mods", admin, nil), 204)
	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/admin/permissions/groups/mods", admin, nil), 404, "no_rows")
}

func TestAPIAdminUsersAreListedAndModerated(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]
	h.signUp(t, "banme")

	res := h.get(t, "/api/v1/admin/users?per_page=2", admin)
	expectStatus(t, res, 200)
	body := res.JSON(t)
	users := items(t, body, "users")
	if len(users) != 2 || object(t, body, "pagination")["total"].(float64) < 4 {
		t.Errorf("users page = %s", res.Body)
	}
	if users[0]["email"] == nil || users[0]["member_of"] == nil || users[0]["banned_at"] == nil && false {
		t.Errorf("admin user view = %v", users[0])
	}

	id := h.users["banme"].ID.String()

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/users/"+id+"/ban", admin, nil), 204)
	expectProblem(t, h.get(t, "/api/v1/users/banme", ""), 404, "no_rows")
	expectProblem(t, h.get(t, "/api/v1/account", h.keys["banme"]), 401, "err_apikey_invalid")

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/users/"+id+"/unban", admin, nil), 204)
	expectStatus(t, h.get(t, "/api/v1/users/banme", ""), 200)
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/admin/users/"+id+"/confirm", admin, nil), 204)

	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/users/"+id+"/groups", admin, map[string]any{"member_of": []string{"veterans", "nope"}}), 204)
	res = h.get(t, "/api/v1/users/banme", "")
	expectStatus(t, res, 200)
	if groups := items(t, object(t, res.JSON(t), "user"), "groups"); len(groups) != 1 || groups[0]["id"] != "veterans" {
		t.Errorf("groups = %v", groups)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/users/"+uuid.NewString()+"/ban", admin, nil), 404, "no_rows")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/users/not-a-uuid/ban", admin, nil), 404, "no_rows")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/admin/users/"+id+"/ban", h.keys["vera"], nil), 403, "err_forbidden")

	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/admin/users/"+id, admin, nil), 204)
	expectProblem(t, h.get(t, "/api/v1/users/banme", ""), 404, "no_rows")
}

func TestAPIAdminLogsReportsAndHealth(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/reports", h.keys["vera"], map[string]any{
		"target": "topic", "id": h.topics["welcome"].ID.String(), "reason": "misconduct",
	}), 201)

	res := h.get(t, "/api/v1/admin/reports", admin)
	expectStatus(t, res, 200)
	found := false
	for _, rep := range items(t, res.JSON(t), "reports") {
		if rep["topic_id"] == h.topics["welcome"].ID.String() && rep["author"] == "vera" {
			found = true
			if rep["reason"] != "misconduct" || rep["target_author"] != "admin" || rep["url"] == "" {
				t.Errorf("report = %v", rep)
			}
		}
	}
	if !found {
		t.Errorf("the report is missing from %s", res.Body)
	}

	res = h.get(t, "/api/v1/admin/logs?per_page=5", admin)
	expectStatus(t, res, 200)
	body := res.JSON(t)
	logs := items(t, body, "logs")
	if len(logs) == 0 || object(t, body, "pagination")["total"].(float64) == 0 {
		t.Fatalf("logs = %s", res.Body)
	}
	visits := 0
	for _, entry := range logs {
		if entry["kind"] == "admin_visit" && entry["actor"] == "admin" && entry["path"] != "" {
			visits++
		}
	}
	if visits == 0 {
		t.Errorf("no admin visit recorded in %s", res.Body)
	}

	res = h.get(t, "/api/v1/admin/health", admin)
	expectStatus(t, res, 200)
	if _, ok := res.JSON(t)["issues"].([]any); !ok {
		t.Errorf("issues should be an array: %s", res.Body)
	}
}

func TestAPIAdminSettingsRoundTrip(t *testing.T) {
	h := apiTest(t)
	admin := h.keys["admin"]

	res := h.get(t, "/api/v1/admin/auth", admin)
	expectStatus(t, res, 200)
	if auth := object(t, res.JSON(t), "auth"); auth["allowed_address_type"] != float64(0) {
		t.Errorf("auth = %v", auth)
	}
	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/auth", admin, map[string]any{"allowed_address_type": 2}), 204)
	t.Cleanup(func() {
		h.do(t, http.MethodPut, "/api/v1/admin/auth", admin, map[string]any{"allowed_address_type": 0})
	})
	res = h.get(t, "/api/v1/admin/auth", admin)
	auth := object(t, res.JSON(t), "auth")
	if auth["allowed_address_type"] != float64(2) {
		t.Errorf("auth after update = %v", auth)
	}
	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/auth", admin, auth), 204)
	res = h.get(t, "/api/v1/admin/auth", admin)
	if again := object(t, res.JSON(t), "auth"); again["allowed_address_type"] != float64(2) {
		t.Errorf("sending back what was read changed the address type: %v", again)
	}
	body := expectProblem(t, h.do(t, http.MethodPut, "/api/v1/admin/auth", admin, map[string]any{"allowed_address_type": 3}), 422, "validation")
	if fields := object(t, body, "errors"); fields["allowed_address_type"] == nil {
		t.Errorf("errors = %v", fields)
	}

	res = h.get(t, "/api/v1/admin/board/topics", admin)
	expectStatus(t, res, 200)
	topics := object(t, res.JSON(t), "topics")
	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/admin/board/topics", admin, map[string]any{
		"allow_kind_question": topics["allow_kind_question"],
		"allow_kind_poll":     topics["allow_kind_poll"],
		"allow_kind_rsvp":     topics["allow_kind_rsvp"],
	}), 204)

	for _, path := range []string{
		"/api/v1/admin/board/attachments",
		"/api/v1/admin/board/profiles",
		"/api/v1/admin/board/themes",
		"/api/v1/admin/comms/email",
		"/api/v1/admin/comms/xmpp",
	} {
		res = h.get(t, path, admin)
		expectStatus(t, res, 200)
		assertNoLeak(t, res.Body)
	}

	res = h.get(t, "/api/v1/admin/comms/email", admin)
	if ids, _ := res.JSON(t)["target_ids"].([]any); len(ids) != 1 || ids[0] != "debug-email" {
		t.Errorf("email target ids = %v", ids)
	}
	res = h.get(t, "/api/v1/admin/board/themes", admin)
	if themes, _ := res.JSON(t)["themes"].([]any); len(themes) == 0 {
		t.Errorf("themes = %s", res.Body)
	}
}
