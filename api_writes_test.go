package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (h *apiHarness) upload(t *testing.T, token string, files map[string]string) apiResponse {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for name, content := range files {
		part, err := w.CreateFormFile("attachments", name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = io.WriteString(part, content); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/attachments", &buf)
	req.Header.Set(fiber.HeaderContentType, w.FormDataContentType())
	if token != "" {
		req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	}

	res, err := h.app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return apiResponse{Status: res.StatusCode, Header: res.Header, Body: raw}
}

func (h *apiHarness) enableAttachments(t *testing.T, enabled bool) {
	t.Helper()

	repo := gh.Repositories(h.rt).Setting
	atts, err := settingRepo.GetByID[setting.Attachments](repo, "attachments")
	if err != nil {
		t.Fatal(err)
	}
	atts.JSONValue.EnableAttachments = enabled
	atts.JSONValue.StorageProviderID = "test-local"
	if err = settingRepo.Update(repo, atts); err != nil {
		t.Fatal(err)
	}
}

func TestAPIWritesNeedAKey(t *testing.T) {
	h := apiTest(t)

	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/topics"},
		{http.MethodPost, "/api/v1/topics/" + h.topics["welcome"].ShortID + "/replies"},
		{http.MethodPost, "/api/v1/topics/" + h.topics["welcome"].ShortID + "/votes"},
		{http.MethodPost, "/api/v1/reports"},
		{http.MethodPost, "/api/v1/attachments"},
	} {
		res := h.do(t, tc.method, tc.path, "", map[string]any{})
		expectProblem(t, res, 401, "err_authentication_required")
		if got := res.Header.Get(fiber.HeaderWWWAuthenticate); got != `Bearer realm="hyperuplink"` {
			t.Errorf("%s %s: WWW-Authenticate = %q", tc.method, tc.path, got)
		}
	}
}

func TestAPITopicCreationFollowsPermissions(t *testing.T) {
	h := apiTest(t)

	res := h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": h.forums["announcements"].String(),
		"name":     "Hello from the API",
		"text":     "Written with **curl**.",
	})
	expectStatus(t, res, 201)
	assertNoLeak(t, res.Body)
	top := object(t, res.JSON(t), "topic")
	if top["slug"] != "hello-from-the-api" || top["kind"] != "regular" || object(t, top, "author")["username"] != "vera" {
		t.Errorf("topic = %v", top)
	}
	if !strings.Contains(top["html"].(string), "<strong>curl</strong>") {
		t.Errorf("html = %v", top["html"])
	}
	if loc := res.Header.Get(fiber.HeaderLocation); loc != "/api/v1/topics/"+top["id"].(string) {
		t.Errorf("Location = %q", loc)
	}
	expectStatus(t, h.get(t, res.Header.Get(fiber.HeaderLocation), ""), 200)

	body := expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": "nope",
		"name":     "",
	}), 422, "validation")
	errs := object(t, body, "errors")
	if errs["forum_id"] == nil || errs["name"] == nil || errs["text"] == nil {
		t.Errorf("errors = %v", errs)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": uuid.NewString(),
		"name":     "Nowhere",
		"text":     "x",
	}), 404, "no_rows")

	h.setDefaultPermission(t, permission.LevelRead)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": h.forums["help"].String(),
		"name":     "Read only",
		"text":     "x",
	}), 403, "err_forbidden")
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["juno"], map[string]any{
		"forum_id": h.forums["help"].String(),
		"name":     "Veterans may post",
		"text":     "x",
	}), 201)

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["vera"], map[string]any{
		"forum_id": h.forums["help"].String(),
		"name":     "Hidden",
		"text":     "x",
	}), 404, "no_rows")
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/topics", h.keys["juno"], map[string]any{
		"forum_id": h.forums["help"].String(),
		"name":     "Veterans still post",
		"text":     "x",
	}), 201)
}

func TestAPIRepliesNestAndNotify(t *testing.T) {
	h := apiTest(t)
	welcome := h.topics["welcome"]

	res := h.do(t, http.MethodPost, "/api/v1/topics/"+welcome.ShortID+"/replies", h.keys["juno"], map[string]any{
		"text":      "A reply from the API.",
		"parent_id": h.replies["first"].String(),
	})
	expectStatus(t, res, 201)
	assertNoLeak(t, res.Body)
	rep := object(t, res.JSON(t), "reply")
	if rep["parent_id"] != h.replies["first"].String() || rep["topic_id"] != welcome.ID.String() {
		t.Errorf("reply = %v", rep)
	}
	if object(t, rep, "author")["username"] != "juno" || !strings.Contains(rep["url"].(string), "#post-"+rep["short_id"].(string)) {
		t.Errorf("reply = %v", rep)
	}
	if loc := res.Header.Get(fiber.HeaderLocation); loc != "/api/v1/replies/"+rep["id"].(string) {
		t.Errorf("Location = %q", loc)
	}
	expectStatus(t, h.get(t, res.Header.Get(fiber.HeaderLocation), ""), 200)

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+welcome.ShortID+"/replies", h.keys["juno"], map[string]any{
		"text":      "Wrong parent",
		"parent_id": uuid.NewString(),
	}), 422, "err_reply_parent_invalid")

	body := expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+welcome.ShortID+"/replies", h.keys["juno"], map[string]any{
		"text": "",
	}), 422, "validation")
	if errs := object(t, body, "errors"); errs["text"] == nil {
		t.Errorf("errors = %v", errs)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/doesnotexist/replies", h.keys["juno"], map[string]any{
		"text": "x",
	}), 404, "no_rows")

	h.setDefaultPermission(t, permission.LevelRead)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+welcome.ShortID+"/replies", h.keys["vera"], map[string]any{
		"text": "Read only",
	}), 403, "err_forbidden")

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+welcome.ShortID+"/replies", h.keys["vera"], map[string]any{
		"text": "Hidden",
	}), 404, "no_rows")
}

func TestAPIPollVotesAreFinal(t *testing.T) {
	h := apiTest(t)

	res := h.do(t, http.MethodPost, "/api/v1/topics", h.keys["admin"], map[string]any{
		"forum_id":     h.forums["announcements"].String(),
		"name":         "Pineapple on pizza?",
		"text":         "Settle it.",
		"kind":         "poll",
		"poll_options": []string{"Yes", "No", "Only on Fridays"},
	})
	expectStatus(t, res, 201)
	top := object(t, res.JSON(t), "topic")
	poll := object(t, top, "poll")
	if top["kind"] != "poll" || len(items(t, poll, "options")) != 3 || poll["ended"] != false {
		t.Fatalf("poll topic = %v", top)
	}
	if viewer := object(t, poll, "viewer"); viewer["can_vote"] != true || viewer["has_voted"] != false {
		t.Errorf("poll viewer after creation = %v", viewer)
	}
	topicID := top["id"].(string)

	res = h.get(t, "/api/v1/topics/"+topicID, "")
	expectStatus(t, res, 200)
	if guestPoll := object(t, object(t, res.JSON(t), "topic"), "poll"); guestPoll["viewer"] != nil {
		t.Errorf("a guest has no poll viewer state: %v", guestPoll)
	}

	res = h.do(t, http.MethodPost, "/api/v1/topics/"+topicID+"/votes", h.keys["vera"], map[string]any{"selection": 2})
	expectStatus(t, res, 201)
	poll = object(t, res.JSON(t), "poll")
	if poll["total"] != float64(1) {
		t.Errorf("total = %v", poll["total"])
	}
	if viewer := object(t, poll, "viewer"); viewer["has_voted"] != true || viewer["selection"] != float64(2) || viewer["can_vote"] != false {
		t.Errorf("viewer = %v", viewer)
	}
	if options := items(t, poll, "options"); options[2]["votes"] != float64(1) || options[2]["percent"] != float64(100) {
		t.Errorf("options = %v", options)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+topicID+"/votes", h.keys["vera"], map[string]any{"selection": 1}), 409, "unique_violation_on")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+topicID+"/votes", h.keys["juno"], map[string]any{"selection": 7}), 422, "err_poll_selection_invalid")
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+h.topics["welcome"].ShortID+"/votes", h.keys["juno"], map[string]any{"selection": 0}), 422, "err_poll_kind_invalid")

	h.setDefaultPermission(t, permission.LevelRead)
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+topicID+"/votes", h.keys["juno"], map[string]any{"selection": 0}), 403, "err_forbidden")
}

func TestAPIReportsArePerViewerAndPost(t *testing.T) {
	h := apiTest(t)

	res := h.do(t, http.MethodPost, "/api/v1/reports", h.keys["vera"], map[string]any{
		"target": "reply",
		"id":     h.replies["second"].String(),
		"reason": "spam",
	})
	expectStatus(t, res, 201)
	if body := res.JSON(t); body["target"] != "reply" || body["reason"] != "spam" {
		t.Errorf("report = %s", res.Body)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/reports", h.keys["vera"], map[string]any{
		"target": "reply",
		"id":     h.replies["second"].String(),
		"reason": "misconduct",
	}), 409, "unique_violation_on")

	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/reports", h.keys["juno"], map[string]any{
		"target": "topic",
		"id":     h.topics["help"].ShortID,
		"reason": "illegal",
	}), 201)

	body := expectProblem(t, h.do(t, http.MethodPost, "/api/v1/reports", h.keys["vera"], map[string]any{
		"target": "comment",
		"id":     "x",
		"reason": "because",
	}), 422, "validation")
	if errs := object(t, body, "errors"); errs["target"] == nil || errs["reason"] == nil {
		t.Errorf("errors = %v", errs)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/reports", h.keys["vera"], map[string]any{
		"target": "topic",
		"id":     uuid.NewString(),
		"reason": "spam",
	}), 404, "no_rows")
}

func TestAPIAttachmentsAreUploadedThenReferenced(t *testing.T) {
	h := apiTest(t)

	h.enableAttachments(t, false)
	expectProblem(t, h.upload(t, h.keys["vera"], map[string]string{"notes.txt": "hello"}), 422, "err_attachments_disabled")

	h.enableAttachments(t, true)
	t.Cleanup(func() { h.enableAttachments(t, false) })

	body := expectProblem(t, h.upload(t, h.keys["vera"], map[string]string{}), 422, "validation")
	if errs := object(t, body, "errors"); errs["attachments"] == nil {
		t.Errorf("errors = %v", errs)
	}

	content := "The modem log, " + time.Now().Format(time.RFC3339Nano) + "\n"
	res := h.upload(t, h.keys["vera"], map[string]string{"modem.log": content})
	expectStatus(t, res, 201)
	atts := items(t, res.JSON(t), "attachments")
	if len(atts) != 1 || atts[0]["filename"] != "modem.log" || atts[0]["mime_type"] != "text/plain; charset=utf-8" {
		t.Fatalf("attachments = %v", atts)
	}
	attID := atts[0]["id"].(string)
	if atts[0]["url"] != "/api/v1/attachments/"+attID {
		t.Errorf("url = %v", atts[0]["url"])
	}

	expectProblem(t, h.upload(t, h.keys["vera"], map[string]string{"again.log": content}), 422, "attachment_duplicate")

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/topics/"+h.topics["help"].ShortID+"/replies", h.keys["juno"], map[string]any{
		"text":           "Not my file",
		"attachment_ids": []string{attID},
	}), 404, "attachment_not_found")

	res = h.do(t, http.MethodPost, "/api/v1/topics/"+h.topics["help"].ShortID+"/replies", h.keys["vera"], map[string]any{
		"text":           "Here is the log.",
		"attachment_ids": []string{attID},
	})
	expectStatus(t, res, 201)
	rep := object(t, res.JSON(t), "reply")
	if got := items(t, rep, "attachments"); len(got) != 1 || got[0]["id"] != attID {
		t.Errorf("reply attachments = %v", got)
	}

	res = h.get(t, "/api/v1/attachments/"+attID, "")
	expectStatus(t, res, 200)
	if string(res.Body) != content || !strings.HasPrefix(res.Header.Get(fiber.HeaderContentType), "text/plain") {
		t.Errorf("download = %q %q", res.Header.Get(fiber.HeaderContentType), res.Body)
	}
	if cd := res.Header.Get(fiber.HeaderContentDisposition); !strings.Contains(cd, `filename="modem.log"`) {
		t.Errorf("Content-Disposition = %q", cd)
	}

	h.setDefaultPermission(t, permission.LevelNone)
	expectProblem(t, h.get(t, "/api/v1/attachments/"+attID, ""), 404, "no_rows")
	expectStatus(t, h.get(t, "/api/v1/attachments/"+attID, h.keys["juno"]), 200)
}
