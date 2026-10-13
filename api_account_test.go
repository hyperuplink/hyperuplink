package main

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/pquerna/otp/totp"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (h *apiHarness) signUp(t *testing.T, name string) string {
	t.Helper()

	usr, err := logicsession.SignUp(h.ctx, h.rt, &logicsession.SignUpInput{
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
			"language":        "en",
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

	res = h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{"language": "de"})
	expectStatus(t, res, 200)
	usr = object(t, res.JSON(t), "user")
	if usr["language"] != "de" || usr["timezone"] != "Europe/Berlin" {
		t.Errorf("language patch = %v", usr)
	}
	langBody := expectProblem(t, h.do(t, http.MethodPatch, "/api/v1/account", h.keys["juno"], map[string]any{"language": "xx"}), 422, "validation")
	if errs := object(t, langBody, "errors"); errs["language"] == nil {
		t.Errorf("errors = %v", errs)
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

func (h *apiHarness) enablePictures(t *testing.T, enabled bool) {
	t.Helper()

	repo := gh.Repositories(h.rt).Setting
	profiles, err := settingRepo.GetByID[setting.Profiles](h.ctx, repo, "profiles")
	if err != nil {
		t.Fatal(err)
	}
	profiles.JSONValue.EnablePicture = enabled
	profiles.JSONValue.PictureFormat = "png"
	profiles.JSONValue.PictureStorageProviderID = "test-local"
	if err = settingRepo.Update(h.ctx, repo, profiles); err != nil {
		t.Fatal(err)
	}
}

func (h *apiHarness) uploadPicture(t *testing.T, token string) apiResponse {
	t.Helper()

	var encoded bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for x := range 4 {
		for y := range 4 {
			img.Set(x, y, color.RGBA{R: uint8(60 * x), G: uint8(60 * y), B: 128, A: 255})
		}
	}
	if err := png.Encode(&encoded, img); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("picture", "me.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(encoded.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/account/picture", &buf)
	req.Header.Set(fiber.HeaderContentType, w.FormDataContentType())
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)

	res, err := h.app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	return apiResponse{Status: res.StatusCode, Header: res.Header, Body: raw}
}

func TestAPIPictureUploadReplacesThePreviousOne(t *testing.T) {
	h := apiTest(t)

	if _, err := exec.LookPath("convert"); err != nil {
		t.Skip("ImageMagick's convert is not installed")
	}

	h.enablePictures(t, true)
	t.Cleanup(func() {
		h.enablePictures(t, false)
		vera, err := h.reloadUser("vera")
		if err != nil {
			t.Error(err)
			return
		}
		vera.ProfilePicture = ""
		if err = gh.Repositories(h.rt).User.Update(h.ctx, vera); err != nil {
			t.Error(err)
		}
		os.RemoveAll(filepath.Join(h.dir, "media", "profile-pictures"))
	})

	pictureFile := func(id string) string {
		return filepath.Join(h.dir, "media", "profile-pictures", id+".png")
	}

	expectStatus(t, h.uploadPicture(t, h.keys["vera"]), 200)
	vera, err := h.reloadUser("vera")
	if err != nil {
		t.Fatal(err)
	}
	first := vera.ProfilePicture
	if first == "" {
		t.Fatal("the first upload stored no picture id")
	}
	if _, err = os.Stat(pictureFile(first)); err != nil {
		t.Fatalf("the first picture isn't in the storage: %v", err)
	}

	expectStatus(t, h.uploadPicture(t, h.keys["vera"]), 200)
	if vera, err = h.reloadUser("vera"); err != nil {
		t.Fatal(err)
	}
	second := vera.ProfilePicture
	if second == "" || second == first {
		t.Fatalf("the second upload stored %q after %q", second, first)
	}
	if _, err = os.Stat(pictureFile(second)); err != nil {
		t.Errorf("the second picture isn't in the storage: %v", err)
	}
	if _, err = os.Stat(pictureFile(first)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the first picture is still in the storage: %v", err)
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
