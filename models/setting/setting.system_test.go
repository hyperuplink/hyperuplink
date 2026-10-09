package setting

import (
	"encoding/json"
	"testing"
)

func TestNewSystemAllowsAPIGuests(t *testing.T) {
	if !NewSystem().AllowAPIGuests {
		t.Error("AllowAPIGuests = false, want true by default")
	}
}

func TestUnmarshalSystemFallsBackToDefaults(t *testing.T) {
	var sys System

	if err := json.Unmarshal([]byte(`{"name":"Board","base_url":"https://b.example","topics_per_page":20}`), &sys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !sys.AllowAPIGuests {
		t.Error("a row stored before allow_api_guests existed decoded it as false, want the default true")
	}
	if sys.Name != "Board" || sys.BaseURL != "https://b.example" || sys.TopicsPerPage != 20 {
		t.Errorf("stored fields were not kept: %+v", sys)
	}
}

func TestUnmarshalSystemKeepsAnExplicitFalse(t *testing.T) {
	var sys System

	if err := json.Unmarshal([]byte(`{"allow_api_guests": false}`), &sys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if sys.AllowAPIGuests {
		t.Error("allow_api_guests: false decoded as true")
	}
}
