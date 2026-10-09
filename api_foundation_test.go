package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/helpers/redisstorage"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func TestAPIDiscoveryIsOpenToEveryone(t *testing.T) {
	h := apiTest(t)

	for name, token := range map[string]string{"guest": "", "user": h.keys["vera"], "admin": h.keys["admin"]} {
		res := h.get(t, "/api/v1", token)
		expectStatus(t, res, 200)

		body := res.JSON(t)
		if body["api"] != "v1" || body["guests"] != true || body["name"] == "" {
			t.Errorf("%s: discovery = %s", name, res.Body)
		}
		links, _ := body["links"].(map[string]any)
		if links["openapi"] != "/api/v1/openapi.json" || links["session"] != "/api/v1/session" {
			t.Errorf("%s: links = %v", name, links)
		}
		if cc := res.Header.Get(fiber.HeaderCacheControl); cc != "private, no-cache" {
			t.Errorf("%s: Cache-Control = %q", name, cc)
		}
		if res.Header.Get(fiber.HeaderETag) == "" {
			t.Errorf("%s: no ETag on a GET", name)
		}
	}
}

func TestAPIUnknownPathIsAProblem(t *testing.T) {
	h := apiTest(t)

	body := expectProblem(t, h.get(t, "/api/v1/nope", ""), 404, "err_not_found")
	if body["instance"] != "/api/v1/nope" {
		t.Errorf("instance = %v", body["instance"])
	}

	expectProblem(t, h.get(t, "/api/v1/nope", h.keys["vera"]), 404, "err_not_found")
	expectProblem(t, h.do(t, http.MethodDelete, "/api/v1", "", nil), 405, "err_method_not_allowed")
}

func TestAPIInvalidKeysAreRejectedWithAChallenge(t *testing.T) {
	h := apiTest(t)

	for _, header := range []string{
		"Bearer hup_doesnotexist",
		"Bearer notaprefixedkey",
		"Basic dmVyYTpodXBlcmh5cGVyIQ==",
		"Bearer",
	} {
		res := h.request(t, h.app, http.MethodGet, "/api/v1", "", nil, map[string]string{
			fiber.HeaderAuthorization: header,
		})
		expectProblem(t, res, 401, "err_apikey_invalid")

		if got := res.Header.Get(fiber.HeaderWWWAuthenticate); got != `Bearer realm="hyperuplink", error="invalid_token"` {
			t.Errorf("%q: WWW-Authenticate = %q", header, got)
		}
	}
}

func TestAPIGuestSwitchClosesEverythingButTheDoor(t *testing.T) {
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

	res := h.get(t, "/api/v1", "")
	expectStatus(t, res, 200)
	if res.JSON(t)["guests"] != false {
		t.Errorf("discovery should report guests: false\n%s", res.Body)
	}

	res = h.get(t, "/api/v1/nope", "")
	expectProblem(t, res, 401, "err_authentication_required")
	if got := res.Header.Get(fiber.HeaderWWWAuthenticate); got != `Bearer realm="hyperuplink"` {
		t.Errorf("WWW-Authenticate = %q", got)
	}

	expectProblem(t, h.get(t, "/api/v1/nope", h.keys["vera"]), 404, "err_not_found")
	expectStatus(t, h.get(t, "/_internal/health/livez", ""), 200)
	expectStatus(t, h.get(t, "/api/v1/openapi.json", ""), 200)
}

func TestAPIRateLimitAnswers429(t *testing.T) {
	h := apiTest(t)

	storage, err := redisstorage.New(h.rt.Config())
	if err != nil {
		t.Fatal(err)
	}
	if err = storage.Delete("ratelimit:ip:0.0.0.0"); err != nil {
		t.Fatal(err)
	}

	_, secret, err := logicapi.Create(h.rt, h.users["vera"].ID, &logicapi.CreateInput{KeyName: "ratelimit"})
	if err != nil {
		t.Fatal(err)
	}

	cfg := h.cfg
	cfg.RateLimitGuest = 2
	cfg.RateLimitUser = 3

	app, err := h.newApp(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for i := range 2 {
		res := h.request(t, app, http.MethodGet, "/api/v1", "", nil, nil)
		if res.Status != 200 {
			t.Fatalf("guest request %d: status %d\n%s", i+1, res.Status, res.Body)
		}
	}
	res := h.request(t, app, http.MethodGet, "/api/v1", "", nil, nil)
	expectProblem(t, res, 429, "err_rate_limited")
	if res.Header.Get(fiber.HeaderRetryAfter) == "" {
		t.Error("no Retry-After on 429")
	}

	for i := range 3 {
		res = h.request(t, app, http.MethodGet, "/api/v1", secret, nil, nil)
		if res.Status != 200 {
			t.Fatalf("user request %d: status %d\n%s", i+1, res.Status, res.Body)
		}
	}
	expectProblem(t, h.request(t, app, http.MethodGet, "/api/v1", secret, nil, nil), 429, "err_rate_limited")

	expectStatus(t, h.request(t, app, http.MethodGet, "/api/v1", h.keys["admin"], nil, nil), 200)
}

func TestAPIHealthAndOpenAPINeedNoKey(t *testing.T) {
	h := apiTest(t)

	expectStatus(t, h.get(t, "/_internal/health/livez", ""), 200)
	expectStatus(t, h.get(t, "/_internal/health/readyz", ""), 200)

	res := h.get(t, "/api/v1/openapi.json", "")
	expectStatus(t, res, 200)
	body := res.JSON(t)
	if body["swagger"] != "2.0" || body["basePath"] != "/api/v1" {
		t.Errorf("openapi document: swagger=%v basePath=%v", body["swagger"], body["basePath"])
	}
	if info, _ := body["info"].(map[string]any); strings.Contains(fmt.Sprint(info["version"]), "{{") {
		t.Errorf("version placeholder not substituted: %v", info["version"])
	}

	res = h.get(t, "/api/v1/docs/index.html", "")
	expectStatus(t, res, 200)
	if ct := res.Header.Get(fiber.HeaderContentType); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("docs content type = %q", ct)
	}
}
