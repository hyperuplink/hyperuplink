package main

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/helpers/redisstorage"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func TestAPISignInIssuesAKeyAndSignOutRevokesIt(t *testing.T) {
	h := apiTest(t)

	res := h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "vera",
		"password": testPassword,
		"name":     "neonmodem on laptop",
	})
	expectStatus(t, res, 201)
	body := res.JSON(t)

	token, _ := body["token"].(string)
	if !strings.HasPrefix(token, "hup_") {
		t.Fatalf("token = %q", token)
	}
	if key := object(t, body, "key"); key["name"] != "neonmodem on laptop" || key["id"] == "" {
		t.Errorf("key = %v", key)
	}
	if usr := object(t, body, "user"); usr["username"] != "vera" || usr["email"] != "vera@example.org" {
		t.Errorf("user = %v", usr)
	}

	res = h.get(t, "/api/v1/session", token)
	expectStatus(t, res, 200)
	if sess := res.JSON(t); sess["authenticated"] != true || object(t, sess, "key")["name"] != "neonmodem on laptop" {
		t.Errorf("session with the new key = %s", res.Body)
	}

	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/session", token, nil), 204)
	expectProblem(t, h.get(t, "/api/v1/session", token), 401, "err_apikey_invalid")

	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/session", "", nil), 401, "err_authentication_required")
}

func TestAPISignInRejectsBadCredentialsAlike(t *testing.T) {
	h := apiTest(t)

	for name, body := range map[string]map[string]any{
		"wrong password": {"username": "vera", "password": "notthepassword"},
		"unknown user":   {"username": "nobody", "password": testPassword},
	} {
		res := h.do(t, http.MethodPost, "/api/v1/session", "", body)
		expectProblem(t, res, 401, "username_password_wrong")
		if got := res.Header.Get("WWW-Authenticate"); got != `Bearer realm="hyperuplink"` {
			t.Errorf("%s: WWW-Authenticate = %q", name, got)
		}
	}

	body := expectProblem(t, h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "v",
	}), 422, "validation")
	errs := object(t, body, "errors")
	if errs["username"] == nil || errs["password"] == nil {
		t.Errorf("errors = %v", errs)
	}

	res := h.request(t, h.app, http.MethodPost, "/api/v1/session", "", nil, map[string]string{
		"Content-Type": "application/json",
	})
	expectProblem(t, res, 400, "err_form_invalid")

	res = h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "vera",
		"password": testPassword,
		"name":     strings.Repeat("x", 65),
	})
	expectProblem(t, res, 422, "validation")
}

func TestAPISignInAsksForTheSecondFactor(t *testing.T) {
	h := apiTest(t)

	otto, err := logicsession.SignUp(h.rt, &logicsession.SignUpInput{
		Username:       "otto",
		Email:          "otto@example.org",
		Password:       testPassword,
		PasswordRepeat: testPassword,
	}, logicsession.SignUpOptions{Activate: true})
	if err != nil {
		t.Fatal(err)
	}

	enrollment, err := user.NewOTPEnrollment("Test", "otto", "")
	if err != nil {
		t.Fatal(err)
	}
	code, err := totp.GenerateCode(enrollment.Secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err = logictwofactor.Enable(h.rt, otto, enrollment.Secret, code); err != nil {
		t.Fatal(err)
	}

	credentials := map[string]any{"username": "otto", "password": testPassword}

	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/session", "", credentials), 401, "err_otp_required")

	credentials["otp_code"] = "000000"
	expectProblem(t, h.do(t, http.MethodPost, "/api/v1/session", "", credentials), 401, "err_otp_code_wrong")

	credentials["otp_code"], err = totp.GenerateCode(enrollment.Secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	res := h.do(t, http.MethodPost, "/api/v1/session", "", credentials)
	expectStatus(t, res, 201)
	if usr := object(t, res.JSON(t), "user"); usr["otp_enabled"] != true {
		t.Errorf("user = %v", usr)
	}
}

func TestAPISignInStaysOpenWhenGuestsAreClosed(t *testing.T) {
	h := apiTest(t)

	repo := gh.Repositories(h.rt).Setting
	sys, err := settingRepo.GetByID[setting.System](repo, "system")
	if err != nil {
		t.Fatal(err)
	}
	sys.JSONValue.AllowAPIGuests = false
	if err = settingRepo.Update(repo, sys); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sys.JSONValue.AllowAPIGuests = true
		if err := settingRepo.Update(repo, sys); err != nil {
			t.Error(err)
		}
	})

	expectProblem(t, h.get(t, "/api/v1/session", ""), 401, "err_authentication_required")

	res := h.do(t, http.MethodPost, "/api/v1/session", "", map[string]any{
		"username": "vera",
		"password": testPassword,
	})
	expectStatus(t, res, 201)
}

func TestAPISignInIsRateLimited(t *testing.T) {
	h := apiTest(t)

	storage, err := redisstorage.New(h.rt.Config())
	if err != nil {
		t.Fatal(err)
	}
	if err = storage.Delete("ratelimit:auth:0.0.0.0:vera"); err != nil {
		t.Fatal(err)
	}

	cfg := h.cfg
	cfg.RateLimitAuth = 2
	app, err := h.newApp(cfg)
	if err != nil {
		t.Fatal(err)
	}

	wrong := map[string]any{"username": "Vera", "password": "notthepassword"}
	for i := range 2 {
		res := h.request(t, app, http.MethodPost, "/api/v1/session", "", wrong, nil)
		if res.Status != 401 {
			t.Fatalf("attempt %d: status %d\n%s", i+1, res.Status, res.Body)
		}
	}
	res := h.request(t, app, http.MethodPost, "/api/v1/session", "", wrong, nil)
	expectProblem(t, res, 429, "err_rate_limited")

	expectStatus(t, h.request(t, app, http.MethodGet, "/api/v1/session", "", nil, nil), 200)
}

func TestAPIKeysAreManagedByTheirOwner(t *testing.T) {
	h := apiTest(t)

	expectProblem(t, h.get(t, "/api/v1/account/keys", ""), 401, "err_authentication_required")

	res := h.get(t, "/api/v1/account/keys", h.keys["juno"])
	expectStatus(t, res, 200)
	assertNoLeak(t, res.Body)
	before := len(items(t, res.JSON(t), "keys"))

	res = h.do(t, http.MethodPost, "/api/v1/account/keys", h.keys["juno"], map[string]any{"name": "laptop"})
	expectStatus(t, res, 201)
	body := res.JSON(t)
	token, _ := body["token"].(string)
	key := object(t, body, "key")
	if !strings.HasPrefix(token, "hup_") || key["name"] != "laptop" {
		t.Fatalf("issued key = %s", res.Body)
	}

	res = h.get(t, "/api/v1/session", token)
	expectStatus(t, res, 200)
	if usr := object(t, res.JSON(t), "user"); usr["username"] != "juno" {
		t.Errorf("the new key acts as %v", usr["username"])
	}

	res = h.get(t, "/api/v1/account/keys", h.keys["juno"])
	expectStatus(t, res, 200)
	if got := len(items(t, res.JSON(t), "keys")); got != before+1 {
		t.Errorf("keys = %d, want %d", got, before+1)
	}

	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/account/keys/"+key["id"].(string), h.keys["vera"], nil), 404, "no_rows")
	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/account/keys/"+uuid.NewString(), h.keys["juno"], nil), 404, "no_rows")
	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1/account/keys/not-a-uuid", h.keys["juno"], nil), 404, "no_rows")

	expectStatus(t, h.do(t, http.MethodDelete, "/api/v1/account/keys/"+key["id"].(string), h.keys["juno"], nil), 204)
	expectProblem(t, h.get(t, "/api/v1/session", token), 401, "err_apikey_invalid")

	body = expectProblem(t, h.do(t, http.MethodPost, "/api/v1/account/keys", h.keys["juno"], map[string]any{}), 422, "validation")
	if errs := object(t, body, "errors"); errs["name"] == nil {
		t.Errorf("errors = %v", errs)
	}
}
