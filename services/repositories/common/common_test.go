package common

import (
	"strings"
	"testing"
)

func TestQueryOrderBy(t *testing.T) {
	tests := []struct {
		name        string
		qo          QueryOptions
		wantContain string
		wantAbsent  string
	}{
		{
			name:        "valid column ascending",
			qo:          QueryOptions{OrderBy: "position", Order: Ascending},
			wantContain: "ORDER BY position ASC",
		},
		{
			name:        "valid qualified column",
			qo:          QueryOptions{OrderBy: "vr.created_at", Order: Descending},
			wantContain: "ORDER BY vr.created_at DESC",
		},
		{
			name:        "empty order defaults to DESC",
			qo:          QueryOptions{OrderBy: "created_at"},
			wantContain: "ORDER BY created_at DESC",
		},
		{
			name:       "injection attempt is dropped",
			qo:         QueryOptions{OrderBy: "id; DROP TABLE users;--"},
			wantAbsent: "ORDER BY",
		},
		{
			name:       "subquery injection is dropped",
			qo:         QueryOptions{OrderBy: "(SELECT password FROM users)"},
			wantAbsent: "ORDER BY",
		},
		{
			name:        "malicious order direction is coerced",
			qo:          QueryOptions{OrderBy: "position", Order: QueryOrder("ASC; DROP TABLE users")},
			wantContain: "ORDER BY position DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.qo.Query("SELECT * FROM t", QueryCapabilities{})
			if tt.wantContain != "" && !strings.Contains(got, tt.wantContain) {
				t.Errorf("query %q missing %q", got, tt.wantContain)
			}
			if tt.wantAbsent != "" && strings.Contains(got, tt.wantAbsent) {
				t.Errorf("query %q must not contain %q", got, tt.wantAbsent)
			}
		})
	}
}

func TestQueryTableIdentifier(t *testing.T) {
	got := QueryOptions{}.Query(
		"SELECT * FROM t",
		QueryCapabilities{Table: "vr", HasDeleted: true},
	)
	if !strings.Contains(got, "vr.deleted_at IS NULL") {
		t.Errorf("expected qualified column, got %q", got)
	}

	got = QueryOptions{}.Query(
		"SELECT * FROM t",
		QueryCapabilities{Table: "vr;DROP TABLE users", HasDeleted: true},
	)
	if strings.Contains(got, "DROP TABLE") {
		t.Errorf("malicious table interpolated: %q", got)
	}
	if !strings.Contains(got, "deleted_at IS NULL") {
		t.Errorf("expected unqualified fallback column, got %q", got)
	}
}
