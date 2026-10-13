package memo

import (
	"errors"
	"testing"
	"time"
)

func counting(calls *int, value int) func() (int, error) {
	return func() (int, error) {
		*calls++
		return value, nil
	}
}

func TestGetLoadsOnceUntilReset(t *testing.T) {
	m := New(time.Minute)
	calls := 0

	for range 3 {
		if got, err := Get(m, "k", counting(&calls, 7)); err != nil || got != 7 {
			t.Fatalf("Get = %d, %v", got, err)
		}
	}
	if calls != 1 {
		t.Errorf("loads = %d, want 1", calls)
	}

	m.Reset()
	if _, err := Get(m, "k", counting(&calls, 7)); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("loads after Reset = %d, want 2", calls)
	}
}

func TestGetReloadsAfterTheTTL(t *testing.T) {
	m := New(20 * time.Millisecond)
	calls := 0

	if _, err := Get(m, "k", counting(&calls, 1)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	if _, err := Get(m, "k", counting(&calls, 1)); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("loads = %d, want 2", calls)
	}
}

func TestGetDoesNotKeepAValueLoadedBeforeAReset(t *testing.T) {
	m := New(time.Minute)
	calls := 0

	if _, err := Get(m, "k", func() (int, error) {
		calls++
		m.Reset()
		return 1, nil
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := Get(m, "k", counting(&calls, 2)); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("loads = %d, want 2: the value loaded before the reset was kept", calls)
	}
}

func TestGetDoesNotCacheAnError(t *testing.T) {
	m := New(time.Minute)
	calls := 0
	boom := errors.New("boom")

	if _, err := Get(m, "k", func() (int, error) { calls++; return 0, boom }); !errors.Is(err, boom) {
		t.Fatalf("err = %v", err)
	}
	if got, err := Get(m, "k", counting(&calls, 3)); err != nil || got != 3 {
		t.Fatalf("Get = %d, %v", got, err)
	}
	if calls != 2 {
		t.Errorf("loads = %d, want 2", calls)
	}
}

func TestGetReloadsOnATypeMismatch(t *testing.T) {
	m := New(time.Minute)
	calls := 0

	if _, err := Get(m, "k", counting(&calls, 1)); err != nil {
		t.Fatal(err)
	}
	if got, err := Get(m, "k", func() (string, error) { calls++; return "s", nil }); err != nil || got != "s" {
		t.Fatalf("Get = %q, %v", got, err)
	}
	if calls != 2 {
		t.Errorf("loads = %d, want 2", calls)
	}
}
