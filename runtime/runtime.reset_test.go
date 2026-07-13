package runtime

import (
	"testing"
	"time"
)

func TestValidateResetTime(t *testing.T) {
	now := time.Date(2026, 7, 13, 10, 42, 30, 0, time.UTC)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"exact match", "10:42", false},
		{"match ignores seconds", "10:42", false},
		{"single digit hour", "9:42", true},
		{"wrong minute", "10:43", true},
		{"wrong hour", "09:42", true},
		{"empty", "", true},
		{"garbage", "not-a-time", true},
		{"seconds included", "10:42:30", true},
		{"hour out of range", "25:42", true},
		{"minute out of range", "10:99", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResetTime(tt.token, now)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateResetTime(%q, %s) error = %v, wantErr %v",
					tt.token, now.Format(ResetTimeLayout), err, tt.wantErr)
			}
		})
	}
}
