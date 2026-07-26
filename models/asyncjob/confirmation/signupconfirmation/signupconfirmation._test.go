package signupconfirmation

import (
	"testing"

	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func TestSetSystemCompletesSignupURL(t *testing.T) {
	entity, err := New(
		&user.User{Username: "alice", Email: "alice@example.org", Language: "en"},
		"Confirm your signup",
		"TOKEN-123456",
		"session/signup",
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if entity.Signup.URL != "" {
		t.Errorf("URL = %q before SetSystem, want empty", entity.Signup.URL)
	}

	entity.SetSystem(&setting.System{
		Name:    "Example Board",
		BaseURL: "https://example.org",
	})

	if got, want := entity.Signup.URL, "https://example.org/session/signup"; got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
	if got, want := entity.System.Name, "Example Board"; got != want {
		t.Errorf("System.Name = %q, want %q", got, want)
	}
}

func TestSetSystemIsIdempotent(t *testing.T) {
	entity, err := New(
		&user.User{Username: "alice", Email: "alice@example.org", Language: "en"},
		"Confirm your signup",
		"TOKEN-123456",
		"session/signup",
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sys := &setting.System{Name: "Example Board", BaseURL: "https://example.org"}

	entity.SetSystem(sys)
	first := entity.Signup.URL

	entity.SetSystem(sys)
	entity.SetSystem(sys)

	if entity.Signup.URL != first {
		t.Errorf("URL = %q after repeated SetSystem, want stable %q",
			entity.Signup.URL, first)
	}
}

func TestSetRecipientCarriesJIDRoutingBit(t *testing.T) {
	entity, err := New(
		&user.User{
			Username:   "bob",
			Email:      "bob@jabber.example.org",
			Language:   "en",
			EmailIsJID: true,
		},
		"Confirm your signup",
		"TOKEN-123456",
		"session/signup",
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if !entity.Recipient.IsJID {
		t.Error("Recipient.IsJID = false, want true for a JID user")
	}
}
