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
	"github.com/pquerna/otp/totp"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
)

func (h *apiHarness) signUp(t *testing.T, name string) string {
	t.Helper()

	usr, err := logicsession.SignUp(h.rt, &logicsession.SignUpInput{
		Username:       name,
		Email:          name + "@example.org",
		Password:       testPassword,
		PasswordRepeat: testPassword,
	}, logicsession.SignUpOptions{Activate: true})
	if err != nil {
		t.Fatal(err)
	}
	h.users[name] = usr

	res := h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": name,
		"password": testPassword,
	})
	expectStatus(t, res, 201)
	token, _ := res.JSON(t)["token"].(string)
	h.keys[name] = token

	return token
}

func TestAPIAccountIsShownAndPatched(t *testing.T) {
	h := apiTest(t)

	expectProblem(t, h.get(t, "/api/v1/account", ""), 401, "err_authentication_required")

	res := h.get(t, "/api/v1/account", h.keys["juno"])
	expectStatus(t, res, 200)
	usr := object(t, res.JSON(t), "user")
	if usr["username"] != "juno" || usr["timezone"] != "UTC" || usr["notify_on_reply"] != true {
		t.Errorf("account = %v", usr)
	}
	if groups := items(t, usr, "groups"); len(groups) != 1 || groups[0]["id"] != "veterans" {
		t.Errorf("groups = %v", groups)
	}

	res = h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{
		"timezone":        "Europe/Berlin",
		"signature_text":  "*73 de juno*",
		"notify_on_reply": false,
	})
	expectStatus(t, res, 200)
	usr = object(t, res.JSON(t), "user")
	if usr["timezone"] != "Europe/Berlin" || usr["notify_on_reply"] != false || !strings.Contains(usr["signature_html"].(string), "<em>73 de juno</em>") {
		t.Errorf("patched account = %v", usr)
	}
	t.Cleanup(func() {
		h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{
			"timezone":        "UTC",
			"signature_text":  "",
			"notify_on_reply": true,
		})
	})

	res = h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{"notify_on_reply": true})
	expectStatus(t, res, 200)
	usr = object(t, res.JSON(t), "user")
	if usr["timezone"] != "Europe/Berlin" || usr["notify_on_reply"] != true {
		t.Errorf("a partial patch changed more than it was given: %v", usr)
	}

	body := expectProblem(t, h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{"timezone": "Nowhere/Land"}), 422, "validation")
	if errs := object(t, body, "errors"); errs["timezone"] == nil {
		t.Errorf("errors = %v", errs)
	}
}

func TestAPIPasswordChangeNeedsTheCurrentOne(t *testing.T) {
	h := apiTest(t)
	token := h.signUp(t, "pwd")

	expectProblem(t, h.do(t, http.MethodPut, "/api/v1/account/password", token, map[string]any{
		"current_password":    "notthepassword",
		"new_password":        "anothersecret1",
		"new_password_repeat": "anothersecret1",
	}), 422, "err_password_wrong")

	body := expectProblem(t, h.do(t, http.MethodPut, "/api/v1/account/password", token, map[string]any{
		"current_password":    testPassword,
		"new_password":        "anothersecret1",
		"new_password_repeat": "mismatch12345",
	}), 422, "validation")
	if errs := object(t, body, "errors"); errs["new_password_repeat"] == nil {
		t.Errorf("errors = %v", errs)
	}

	expectStatus(t, h.do(t, http.MethodPut, "/api/v1/account/password", token, map[string]any{
		"current_password":    testPassword,
		"new_password":        "anothersecret1",
		"new_password_repeat": "anothersecret1",
	}), 204)

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "pwd", "password": testPassword,
	}), 401, "username_password_wrong")
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "pwd", "password": "anothersecret1",
	}), 201)
	expectStatus(t, h.get(t, "/api/v1/account", token), 200)
}

func TestAPITwoFactorIsEnrolledEnabledAndDisabled(t *testing.T) {
	h := apiTest(t)
	token := h.signUp(t, "tfa")

	res := h.get(t, "/api/v1/account/twofactor", token)
	expectStatus(t, res, 200)
	body := res.JSON(t)
	enrollment := object(t, body, "enrollment")
	secret, _ := enrollment["secret"].(string)
	if body["enabled"] != false || secret == "" || !strings.HasPrefix(enrollment["url"].(string), "otpauth://") || !strings.HasPrefix(enrollment["qr_data_uri"].(string), "data:image/png;base64,") {
		t.Fatalf("enrollment = %s", res.Body)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/account/twofactor", token, map[string]any{
		"otp_url": "not-a-url", "otp_code": "123456",
	}), 422, "validation")

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/account/twofactor", token, map[string]any{
		"otp_url": enrollment["url"], "otp_code": "000000",
	}), 422, "err_otp_code_wrong")

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	expectStatus(t, h.do(t, http.MethodPost, "/api/v1/account/twofactor", token, map[string]any{
		"otp_url": enrollment["url"], "otp_code": code,
	}), 204)

	res = h.get(t, "/api/v1/account/twofactor", token)
	expectStatus(t, res, 200)
	if body = res.JSON(t); body["enabled"] != true || body["enrollment"] != nil {
		t.Errorf("after enabling = %s", res.Body)
	}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "tfa", "password": testPassword,
	}), 401, "err_otp_required")

	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/account/twofactor", token, map[string]any{
		"current_password": "notthepassword",
	}), 422, "err_password_wrong")
	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/account/twofactor", token, map[string]any{
		"current_password": testPassword,
	}), 204)

	res = h.get(t, "/api/v1/account/twofactor", token)
	expectStatus(t, res, 200)
	if res.JSON(t)["enabled"] != false {
		t.Errorf("after disabling = %s", res.Body)
	}
}

func TestAPIPictureUploadIsRejectedWhileDisabled(t *testing.T) {
	h := apiTest(t)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("picture", "me.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = io.WriteString(part, "not really a png"); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/account/picture", &buf)
	req.Header.Set(fiber.HeaderContentType, w.FormDataContentType())
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+h.keys["vera"])

	res, err := h.app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	expectProblem(t, apiResponse{Status: res.StatusCode, Header: res.Header, Body: raw}, 422, "err_pictures_disabled")
}
