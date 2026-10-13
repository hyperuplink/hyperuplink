package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
)

type webClient struct {
	t       *testing.T
	app     *fiber.App
	lang    string
	cookies map[string]string
}

func newWebClient(t *testing.T, app *fiber.App, lang string) *webClient {
	return &webClient{t: t, app: app, lang: lang, cookies: map[string]string{}}
}

func (wc *webClient) do(method string, path string, form url.Values) *http.Response {
	wc.t.Helper()

	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req := httptest.NewRequest(method, path, body)
	if form != nil {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	}
	if wc.lang != "" {
		req.Header.Set(fiber.HeaderAcceptLanguage, wc.lang)
	}
	for name, value := range wc.cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	res, err := wc.app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		wc.t.Fatalf("%s %s: %v", method, path, err)
	}
	for _, cookie := range res.Cookies() {
		wc.cookies[cookie.Name] = cookie.Value
	}

	return res
}

func (wc *webClient) get(path string) (int, string) {
	wc.t.Helper()

	res := wc.do(http.MethodGet, path, nil)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		wc.t.Fatalf("GET %s: read body: %v", path, err)
	}

	return res.StatusCode, string(body)
}

var csrfField = regexp.MustCompile(`name="_csrf" value="([^"]+)"`)

func (wc *webClient) signIn(username string, password string) {
	wc.t.Helper()

	status, body := wc.get("/session/signin")
	if status != 200 {
		wc.t.Fatalf("sign-in page = %d", status)
	}
	token := csrfField.FindStringSubmatch(body)
	if token == nil {
		wc.t.Fatal("the sign-in page has no CSRF token")
	}

	res := wc.do(http.MethodPost, "/session/signin", url.Values{
		"username": {username},
		"password": {password},
		"_csrf":    {token[1]},
	})
	res.Body.Close()
	if res.StatusCode != 302 && res.StatusCode != 303 {
		wc.t.Fatalf("sign-in = %d, want a redirect", res.StatusCode)
	}
}

func TestWebGuestsGetTheLanguageOfTheirBrowser(t *testing.T) {
	_, app := webTest(t)

	for header, want := range map[string]string{
		"de-CH,de;q=0.9,en;q=0.5": "Anmelden",
		"pt-BR,pt;q=0.9":          "Sign in",
		"":                        "Sign in",
	} {
		status, body := newWebClient(t, app, header).get("/session/signin")
		if status != 200 {
			t.Fatalf("Accept-Language %q: status %d", header, status)
		}
		if !strings.Contains(body, want) {
			t.Errorf("Accept-Language %q: the sign-in page doesn't show %q", header, want)
		}
	}
}

func TestWebSignedInUsersGetTheirStoredLanguage(t *testing.T) {
	h, app := webTest(t)
	juno := h.users["juno"]

	if err := logicsettings.Update(h.ctx, h.rt, juno, &logicsettings.UpdateInput{
		Language: "de",
		Timezone: juno.Timezone,
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := logicsettings.Update(h.ctx, h.rt, juno, &logicsettings.UpdateInput{
			Language: "en",
			Timezone: juno.Timezone,
		}); err != nil {
			t.Errorf("restore the language: %v", err)
		}
	})

	wc := newWebClient(t, app, "en-US,en;q=0.9")
	wc.signIn("juno", testPassword)

	status, body := wc.get("/account/settings")
	if status != 200 {
		t.Fatalf("settings page = %d", status)
	}
	if !strings.Contains(body, "Einstellungen") {
		t.Error("the settings page of a user with language de isn't in German")
	}
	if !regexp.MustCompile(`value="de"\s+selected`).MatchString(body) {
		t.Error("the language select doesn't preselect de")
	}
}

func TestSignUpStoresASupportedLanguage(t *testing.T) {
	h := apiTest(t)

	for i, tc := range []struct{ given, want string }{
		{"pt-BR", "en"},
		{"de-AT", "de"},
		{"", "en"},
		{"xx", "en"},
	} {
		usr, err := logicsession.SignUp(h.ctx, h.rt, &logicsession.SignUpInput{
			Username:       fmt.Sprintf("lang%d", i),
			Email:          fmt.Sprintf("lang%d@example.org", i),
			Password:       testPassword,
			PasswordRepeat: testPassword,
			Language:       tc.given,
		}, logicsession.SignUpOptions{Activate: true})
		if err != nil {
			t.Fatalf("signup %d: %v", i, err)
		}
		if usr.Language != tc.want {
			t.Errorf("language given %q stored as %q, want %q", tc.given, usr.Language, tc.want)
		}
	}
}
