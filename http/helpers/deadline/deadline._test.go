package deadline

import (
	"context"
	"errors"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/errs"
	"xn--gckvb8fzb.com/glides/services/database"
)

func TestNewSetsTheDeadlineOfTheRequestContext(t *testing.T) {
	app := fiber.New()
	app.Use(New(time.Second))

	var deadline time.Time
	var ok bool
	app.Get("/", func(c fiber.Ctx) error {
		deadline, ok = c.Context().Deadline()
		return nil
	})

	start := time.Now()
	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	if err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("the request = %v, %v", resp, err)
	}
	if !ok || deadline.Before(start.Add(time.Second)) || deadline.After(time.Now().Add(time.Second)) {
		t.Errorf("the deadline = %v, %v, want one second after %v", deadline, ok, start)
	}
}

func TestNewWithoutATimeoutSetsNoDeadline(t *testing.T) {
	app := fiber.New()
	app.Use(New(0))

	ok := true
	app.Get("/", func(c fiber.Ctx) error {
		_, ok = c.Context().Deadline()
		return nil
	})

	if _, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil)); err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("a timeout of zero set a deadline")
	}
}

func TestNewCancelsAStatementAtTheDeadline(t *testing.T) {
	connection := os.Getenv("HUP_TEST_DATABASE")
	if connection == "" {
		t.Skip("set HUP_TEST_DATABASE to a PostgreSQL URL")
	}

	db, err := database.New(nil, connection)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Startup(); err != nil {
		t.Fatal(err)
	}
	defer db.Shutdown()

	app := fiber.New()
	app.Use(New(50 * time.Millisecond))

	var serr error
	app.Get("/", func(c fiber.Ctx) error {
		_, serr = database.Exec(c.Context(), db, `SELECT pg_sleep(2)`)
		return nil
	})

	start := time.Now()
	if _, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil), fiber.TestConfig{Timeout: 5 * time.Second}); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(serr, errs.ErrUnavailable) || !errors.Is(serr, context.DeadlineExceeded) {
		t.Errorf("the statement ended with %v, want the deadline inside errs.ErrUnavailable", serr)
	}
	if time.Since(start) > time.Second {
		t.Errorf("the statement ran for %v past a deadline of 50ms", time.Since(start))
	}
}
