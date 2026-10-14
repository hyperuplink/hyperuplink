package errs

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestPublicKeepsWhatAVisitorMaySee(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{"sentinel", ErrPictureTooLarge, ErrPictureTooLarge.Error()},
		{"wrapped sentinel", fmt.Errorf("store: %w", ErrAttachmentDuplicate), ErrAttachmentDuplicate.Error()},
		{"validation key", fmt.Errorf("%w_email_jid", ErrValidation), "validation_email_jid"},
		{"unique violation with its constraint", fmt.Errorf("%w_%s", ErrUniqueViolationOn, "idx_users_username"), "unique_violation_on"},
		{"unavailable with the driver's text", fmt.Errorf("%w: %w", ErrUnavailable, errors.New("dial tcp 127.0.0.1:5432: connect: connection refused")), "err_unavailable"},
		{"database error", errors.New(`ERROR: new row for relation "users" violates check constraint "x" (SQLSTATE 23514)`), "err_internal"},
		{"internal", ErrInternal, "err_internal"},
		{"internal sentinel", ErrRedisAddrsEmpty, "err_internal"},
		{"validation itself", ErrValidation, "validation"},
	} {
		if got := Public(tc.err).Error(); got != tc.want {
			t.Errorf("%s: Public = %q, want %q", tc.name, got, tc.want)
		}
	}

	if Public(nil) != nil {
		t.Error("Public(nil) is not nil")
	}
}

func TestEveryPublicErrorHasAnEnglishMessage(t *testing.T) {
	catalog, err := os.ReadFile("../locales/en.toml")
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?m)^([A-Za-z0-9_]+)\s*=`).FindAllSubmatch(catalog, -1) {
		keys[string(m[1])] = true
	}

	if len(public) < 4 {
		t.Fatalf("only %d public errors are registered", len(public))
	}
	for _, err := range public {
		if !keys[err.Error()] {
			t.Errorf("%q is public but locales/en.toml has no message for it", err)
		}
	}
}
