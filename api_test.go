package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api"
	"xn--gckvb8fzb.com/hyperuplink/http/routes"
	logicperms "xn--gckvb8fzb.com/hyperuplink/logic/helpers/perms"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

const (
	testPassword = "hyperhyper!"
	testTimeout  = 10 * time.Second
)

type apiHarness struct {
	ctx        context.Context
	rt         *runtime.Runtime
	cfg        gh.APIConfig
	app        *fiber.App
	dir        string
	users      map[string]*user.User
	keys       map[string]string
	categories map[string]uuid.UUID
	forums     map[string]uuid.UUID
	topics     map[string]*vtopic.VTopic
	replies    map[string]uuid.UUID
}

type apiResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

func (res apiResponse) JSON(t *testing.T) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(res.Body, &body); err != nil {
		t.Fatalf("response is not a JSON object: %v\n%s", err, res.Body)
	}

	return body
}

func (res apiResponse) Code(t *testing.T) string {
	t.Helper()

	code, _ := res.JSON(t)["code"].(string)

	return code
}

var (
	harnessOnce sync.Once
	harness     *apiHarness
	harnessErr  error
)

func apiTest(t *testing.T) *apiHarness {
	t.Helper()

	dsn := os.Getenv("HUP_TEST_DATABASE")
	if dsn == "" {
		t.Skip("set HUP_TEST_DATABASE to a PostgreSQL URL the tests may wipe")
	}

	harnessOnce.Do(func() {
		harness, harnessErr = newAPIHarness(dsn)
	})
	if harnessErr != nil {
		t.Fatalf("api harness: %v", harnessErr)
	}

	return harness
}

func newAPIHarness(dsn string) (h *apiHarness, err error) {
	h = &apiHarness{
		ctx:        context.Background(),
		users:      map[string]*user.User{},
		keys:       map[string]string{},
		categories: map[string]uuid.UUID{},
		forums:     map[string]uuid.UUID{},
		topics:     map[string]*vtopic.VTopic{},
		replies:    map[string]uuid.UUID{},
	}

	if h.dir, err = os.MkdirTemp("", "hyperuplink-api-test-"); err != nil {
		return nil, err
	}

	cfgPath := filepath.Join(h.dir, "hyperuplink.toml")
	if err = os.WriteFile(cfgPath, []byte(h.config(dsn)), 0o600); err != nil {
		return nil, err
	}

	if err = resetSchema(dsn); err != nil {
		return nil, fmt.Errorf("reset schema: %w", err)
	}

	if h.rt, err = newRuntime("file://" + cfgPath); err != nil {
		return nil, err
	}

	h.rt.AddEmbed("migrations", &embedMigrations)
	database.From(h.rt).SetMigrations(h.rt.GetEmbed("migrations"))

	if err = prepare(h.rt); err != nil {
		return nil, err
	}
	if err = h.rt.Startup(); err != nil {
		return nil, err
	}

	routes.Use()

	if err = h.seed(); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	if h.cfg, err = gh.APIConfigFrom(h.rt.Config()); err != nil {
		return nil, err
	}
	if h.app, err = h.newApp(h.cfg); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *apiHarness) config(dsn string) string {
	return fmt.Sprintf(`[General]
Mode = "production"

[Logging]
Level = "error"

[Redis]
Addresses = ["127.0.0.1:6379"]
Database = 15

[Database]
Connection = %q

[Web]
Enable = false

[API]
Enable = true
RateLimitGuest = 0
RateLimitUser = 0
RateLimitAuth = 0

[Users]
PromoteAdmin = ["admin@example.org"]

[[Target]]
ID = "debug-email"
Type = "debug"
Debug.Emulates = "email"
Debug.Path = %q

[[Storage]]
ID = "test-local"
Type = "Local"
Local.Path = %q
Local.PublicURI = "/media"
`, dsn, filepath.Join(h.dir, "messages"), filepath.Join(h.dir, "media"))
}

func resetSchema(dsn string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public")

	return err
}

func (h *apiHarness) newApp(cfg gh.APIConfig) (*fiber.App, error) {
	srv, err := api.New(h.rt, cfg)
	if err != nil {
		return nil, err
	}
	if err = srv.Startup(); err != nil {
		return nil, err
	}

	return srv.App(), nil
}

func (h *apiHarness) perms(usr *user.User) *permission.Resolution {
	return logicperms.Resolve(h.ctx, h.rt, usr.Role, usr.MemberOf)
}

func (h *apiHarness) seed() (err error) {
	for _, u := range []struct{ name, email string }{
		{"admin", "admin@example.org"},
		{"vera", "vera@example.org"},
		{"juno", "juno@example.org"},
	} {
		usr, serr := logicsession.SignUp(h.ctx, h.rt, &logicsession.SignUpInput{
			Username:       u.name,
			Email:          u.email,
			Password:       testPassword,
			PasswordRepeat: testPassword,
		}, logicsession.SignUpOptions{Activate: true})
		if serr != nil {
			return fmt.Errorf("user %s: %w", u.name, serr)
		}
		h.users[u.name] = usr

		_, secret, kerr := logicapi.Create(h.ctx, h.rt, usr.ID, &logicapi.CreateInput{KeyName: "test"})
		if kerr != nil {
			return kerr
		}
		h.keys[u.name] = secret
	}

	for _, c := range []struct{ name, slug string }{
		{"General", "general"},
		{"Support", "support"},
	} {
		id, cerr := logiccategories.Create(h.ctx, h.rt, &logiccategories.CreateInput{Name: c.name, Slug: c.slug})
		if cerr != nil {
			return cerr
		}
		h.categories[c.slug] = id
	}

	for _, f := range []struct{ name, slug, category string }{
		{"Announcements", "announcements", "general"},
		{"Help", "help", "support"},
	} {
		id, ferr := logicforums.Create(h.ctx, h.rt, &logicforums.CreateInput{
			Name:       f.name,
			Slug:       f.slug,
			CategoryID: h.categories[f.category].String(),
		})
		if ferr != nil {
			return ferr
		}
		h.forums[f.slug] = id
	}

	if err = logicpermissions.GroupCreate(h.ctx, h.rt, &logicpermissions.GroupCreateInput{
		ID:   "veterans",
		Name: "Veterans",
	}); err != nil {
		return err
	}
	if err = logicpermissions.Apply(h.ctx, h.rt, &logicpermissions.ApplyInput{
		GroupID:    "veterans",
		CategoryID: h.categories["support"].String(),
		Level:      permission.LevelReadWrite,
	}); err != nil {
		return err
	}
	if err = logicuser.UpdateMembership(h.ctx, h.rt, "juno", &logicuser.MembershipInput{
		MemberOf: []string{"veterans"},
	}); err != nil {
		return err
	}
	if h.users["juno"], err = h.reloadUser("juno"); err != nil {
		return err
	}

	admin := h.users["admin"]
	vera := h.users["vera"]

	welcome, err := logicnewpost.Create(h.ctx, h.rt, admin.ID, h.perms(admin), time.UTC, &logicnewpost.CreateInput{
		Name:    "Welcome to the board",
		Text:    "Be excellent to each other.",
		ForumID: h.forums["announcements"].String(),
		Kind:    "regular",
	}, nil)
	if err != nil {
		return err
	}
	h.topics["welcome"] = welcome

	help, err := logicnewpost.Create(h.ctx, h.rt, vera.ID, h.perms(vera), time.UTC, &logicnewpost.CreateInput{
		Name:    "Need help with my modem",
		Text:    "It only connects at 2400 baud.",
		ForumID: h.forums["help"].String(),
		Kind:    "regular",
	}, nil)
	if err != nil {
		return err
	}
	h.topics["help"] = help

	first, err := logictopics.CreateReply(h.ctx, h.rt, vera.ID, h.perms(vera), &logictopics.CreateReplyInput{
		Text:    "Thanks, happy to be here.",
		TopicID: welcome.ID.String(),
	}, nil)
	if err != nil {
		return err
	}
	h.replies["first"] = first.Reply.ID

	second, err := logictopics.CreateReply(h.ctx, h.rt, admin.ID, h.perms(admin), &logictopics.CreateReplyInput{
		Text:    "Welcome, vera.",
		TopicID: welcome.ID.String(),
		ReplyID: first.Reply.ID.String(),
	}, nil)
	if err != nil {
		return err
	}
	h.replies["second"] = second.Reply.ID

	return nil
}

func queryOne() common.QueryOptions {
	return common.QueryOptions{Limit: 1}
}

func (h *apiHarness) reloadUser(username string) (*user.User, error) {
	return gh.Repositories(h.rt).User.GetByUsername(h.ctx, username, queryOne())
}

func (h *apiHarness) setDefaultPermission(t *testing.T, level string) {
	t.Helper()

	if err := logicpermissions.Apply(h.ctx, h.rt, &logicpermissions.ApplyInput{Level: level}); err != nil {
		t.Fatalf("default permission %s: %v", level, err)
	}
	t.Cleanup(func() {
		if err := logicpermissions.Apply(h.ctx, h.rt, &logicpermissions.ApplyInput{Level: permission.LevelReadWrite}); err != nil {
			t.Errorf("restore default permission: %v", err)
		}
	})
}

func (h *apiHarness) request(
	t *testing.T,
	app *fiber.App,
	method string,
	path string,
	token string,
	body any,
	headers map[string]string,
) apiResponse {
	t.Helper()

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	}
	if token != "" {
		req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := app.Test(req, fiber.TestConfig{Timeout: testTimeout, FailOnTimeout: true})
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return apiResponse{Status: res.StatusCode, Header: res.Header, Body: raw}
}

func (h *apiHarness) get(t *testing.T, path string, token string) apiResponse {
	t.Helper()
	return h.request(t, h.app, http.MethodGet, path, token, nil, nil)
}

func (h *apiHarness) do(t *testing.T, method string, path string, token string, body any) apiResponse {
	t.Helper()
	return h.request(t, h.app, method, path, token, body, nil)
}

func expectStatus(t *testing.T, res apiResponse, status int) {
	t.Helper()

	if res.Status != status {
		t.Fatalf("status = %d, want %d\n%s", res.Status, status, res.Body)
	}
}

func expectProblem(t *testing.T, res apiResponse, status int, code string) map[string]any {
	t.Helper()

	expectStatus(t, res, status)
	if ct := res.Header.Get(fiber.HeaderContentType); ct != "application/problem+json" {
		t.Errorf("content type = %q, want application/problem+json", ct)
	}

	body := res.JSON(t)
	if body["code"] != code {
		t.Errorf("code = %v, want %q\n%s", body["code"], code, res.Body)
	}
	if body["status"] != float64(status) {
		t.Errorf("status member = %v, want %d", body["status"], status)
	}
	if body["type"] != "about:blank" || body["title"] == "" {
		t.Errorf("type = %v, title = %v", body["type"], body["title"])
	}

	return body
}
