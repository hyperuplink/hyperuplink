package setting

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestUserProfileIDIsPrefixedUserUUID(t *testing.T) {
	id := uuid.MustParse("733fc4b0-a296-4d2d-8ed5-c69b3734d6bc")

	got := UserProfileID(id)
	want := "profile-733fc4b0-a296-4d2d-8ed5-c69b3734d6bc"

	if got != want {
		t.Errorf("UserProfileID() = %q, want %q", got, want)
	}

	if len(got) > 64 {
		t.Errorf("UserProfileID() is %d chars, settings.id is VARCHAR(64)", len(got))
	}
}

func TestNewUserProfileEnablesReplyNotifications(t *testing.T) {
	if !NewUserProfile().NotifyOnReply {
		t.Error("NotifyOnReply = false, want true by default")
	}
}

func TestNewUserProfileEnablesEveryViewToggle(t *testing.T) {
	profile := NewUserProfile()

	for name, got := range map[string]bool{
		"ShowBanner":          profile.ShowBanner,
		"ShowFooter":          profile.ShowFooter,
		"ShowProfilePictures": profile.ShowProfilePictures,
	} {
		if !got {
			t.Errorf("%s = false, want true by default", name)
		}
	}
}

func TestUnmarshalUserProfileFallsBackToDefaults(t *testing.T) {
	var profile UserProfile

	if err := json.Unmarshal([]byte(`{"notify_on_reply": true}`), &profile); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	want := NewUserProfile()
	if profile != want {
		t.Errorf("decoding a row stored before the view toggles existed gave "+
			"%+v, want %+v; a key missing from the stored JSON must fall back "+
			"to the Go-side default instead of decoding as false, or every "+
			"user predating the field silently gets the setting turned off",
			profile, want)
	}
}

func TestUnmarshalUserProfileKeepsStoredValues(t *testing.T) {
	var profile UserProfile

	if err := json.Unmarshal(
		[]byte(`{"show_banner": false, "show_footer": false}`),
		&profile,
	); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if profile.ShowBanner || profile.ShowFooter {
		t.Errorf("ShowBanner = %t, ShowFooter = %t, want both false; "+
			"a stored value must override the default",
			profile.ShowBanner, profile.ShowFooter)
	}
	if !profile.ShowProfilePictures {
		t.Error("ShowProfilePictures = false, want the default true to survive")
	}
}

func TestUserProfileViewNamesMatchTheMarshalledFields(t *testing.T) {
	raw, err := json.Marshal(NewUserProfile())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, view := range []string{
		UserProfileViewBanner,
		UserProfileViewFooter,
		UserProfileViewProfilePictures,
	} {
		if _, ok := decoded["show_"+view]; !ok {
			t.Errorf("view %q has no matching show_%s key (got %v); the View "+
				"menu links to /account/settings/view/%s, so the two must "+
				"stay in step or the toggle stops resolving to a field",
				view, view, decoded, view)
		}
	}
}

func TestNotifyOnReplyKeyMatchesTheMarshalledField(t *testing.T) {
	raw, err := json.Marshal(NewUserProfile())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := decoded[UserProfileNotifyOnReplyKey]; !ok {
		t.Errorf("marshalled UserProfile has no %q key (got %v); "+
			"the recipient query filters on this key, so a struct tag "+
			"rename would silently stop honouring the setting",
			UserProfileNotifyOnReplyKey, decoded)
	}
}
