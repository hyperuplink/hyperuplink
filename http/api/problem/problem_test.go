package problem

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
)

func TestFromMapsKnownErrorsAndHidesUnknownOnes(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
		code   string
	}{
		{errs.ErrNoRows, 404, "no_rows"},
		{errs.ErrForbidden, 403, "err_forbidden"},
		{errs.ErrAPIKeyInvalid, 401, "err_apikey_invalid"},
		{errs.ErrAuthenticationRequired, 401, "err_authentication_required"},
		{errs.ErrUniqueViolationOn, 409, "unique_violation_on"},
		{errs.ErrRateLimited, 429, "err_rate_limited"},
		{errs.ErrUnavailable, 503, "err_unavailable"},
		{errs.ErrFormInvalid, 400, "err_form_invalid"},
		{errs.ErrPollEnded, 422, "err_poll_ended"},
		{fmt.Errorf("%w_email_jid", errs.ErrValidation), 422, "validation_email_jid"},
		{fmt.Errorf("%w: %w", errs.ErrUnavailable, errors.New("dial tcp: refused")), 503, "err_unavailable"},
		{fmt.Errorf("%w_idx_postevents_type", errs.ErrUniqueViolationOn), 409, "unique_violation_on"},
		{errors.New("ERROR: relation \"users\" does not exist (SQLSTATE 42P01)"), 500, "err_internal"},
	} {
		p := From(tc.err)
		if p.Status != tc.status || p.Code != tc.code {
			t.Errorf("From(%v) = %d %q, want %d %q", tc.err, p.Status, p.Code, tc.status, tc.code)
		}
		if p.Type != TypeBlank || p.Title == "" {
			t.Errorf("From(%v) has type %q and title %q", tc.err, p.Type, p.Title)
		}
	}
}

func TestFromFiberUsesTheStatusText(t *testing.T) {
	for _, tc := range []struct {
		err    error
		status int
		code   string
	}{
		{fiber.ErrNotFound, 404, "err_not_found"},
		{fiber.ErrMethodNotAllowed, 405, "err_method_not_allowed"},
		{fiber.ErrRequestEntityTooLarge, 413, "err_request_entity_too_large"},
		{fiber.NewError(fiber.StatusTeapot), 418, "err_im_a_teapot"},
		{errors.New("runtime error: index out of range"), 500, "err_internal"},
		{errs.ErrNoRows, 404, "no_rows"},
	} {
		p := FromFiber(tc.err)
		if p.Status != tc.status || p.Code != tc.code {
			t.Errorf("FromFiber(%v) = %d %q, want %d %q", tc.err, p.Status, p.Code, tc.status, tc.code)
		}
	}
}

func TestValidationCarriesTheFields(t *testing.T) {
	p := Validation(map[string][]string{"forum_id": {"uuid"}, "name": {"required"}})

	if p.Status != 422 || p.Code != "validation" {
		t.Errorf("Validation() = %d %q", p.Status, p.Code)
	}
	if len(p.Errors) != 2 || p.Errors["forum_id"][0] != "uuid" {
		t.Errorf("Validation() errors = %v", p.Errors)
	}
}

func TestChallengeNamesInvalidTokensOnly(t *testing.T) {
	if got := Challenge(From(errs.ErrAPIKeyInvalid)); got != `Bearer realm="hyperuplink", error="invalid_token"` {
		t.Errorf("invalid key challenge = %q", got)
	}
	if got := Challenge(From(errs.ErrAuthenticationRequired)); got != `Bearer realm="hyperuplink"` {
		t.Errorf("missing key challenge = %q", got)
	}
}
