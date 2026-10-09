package view

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestTimeMarshalsAsUTCWithAZone(t *testing.T) {
	ts := pgtype.Timestamp{
		Time:  time.Date(2026, 7, 20, 9, 13, 49, 793509000, time.UTC),
		Valid: true,
	}

	got, err := json.Marshal(NewTime(ts))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `"2026-07-20T09:13:49.793509Z"` {
		t.Errorf("Time = %s", got)
	}

	got, err = json.Marshal(NewTime(pgtype.Timestamp{}))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "null" {
		t.Errorf("invalid Time = %s, want null", got)
	}
}

func TestPaginationCountsPages(t *testing.T) {
	p := NewPagination(2, 10, 12)
	if p.Pages != 2 || p.Page != 2 || p.PerPage != 10 || p.Total != 12 {
		t.Errorf("NewPagination(2, 10, 12) = %+v", p)
	}
	if NewPagination(1, 10, 0).Pages != 0 {
		t.Error("an empty collection should have no pages")
	}
}
