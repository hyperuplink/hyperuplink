package permission

import (
	"testing"

	"github.com/google/uuid"
)

func TestCapLowersEveryLevelAboveTheCap(t *testing.T) {
	open := uuid.New()
	closed := uuid.New()

	res := NewResolution()
	res.Default = ReadWrite
	res.SetCategory(open, "open", Moderate)
	res.SetCategory(closed, "closed", None)

	res.Cap(ReadOnly)

	if res.Default != ReadOnly {
		t.Errorf("Default = %03b, want %03b", res.Default, ReadOnly)
	}
	if res.LevelByID(open) != ReadOnly || res.LevelBySlug("open") != ReadOnly {
		t.Errorf("open category = %03b/%03b, want read", res.LevelByID(open), res.LevelBySlug("open"))
	}
	if res.LevelByID(closed) != None {
		t.Errorf("closed category = %03b, want none to stay none", res.LevelByID(closed))
	}
	if res.CanWriteAny() {
		t.Error("CanWriteAny() = true after a cap at read")
	}
	if !res.CanReadID(open) || res.CanReadID(closed) {
		t.Error("reading changed under the cap")
	}
}

func TestLevelsReportsEveryCategory(t *testing.T) {
	a, b := uuid.New(), uuid.New()

	res := NewResolution()
	res.Default = ReadOnly
	res.SetCategory(a, "a", ReadWrite)
	res.SetCategory(b, "b", ReadOnly)

	levels := res.Levels()
	if len(levels) != 2 || levels[a] != ReadWrite || levels[b] != ReadOnly {
		t.Errorf("Levels() = %v", levels)
	}
	if res.DefaultLevel() != ReadOnly {
		t.Errorf("DefaultLevel() = %03b, want read", res.DefaultLevel())
	}

	admin := &Resolution{IsAdmin: true}
	if admin.DefaultLevel() != Moderate || len(admin.Levels()) != 0 {
		t.Errorf("an administrator should report moderate as default and no per-category levels")
	}
}

func TestReadAllNoticesAClosedCategory(t *testing.T) {
	open, closed := uuid.New(), uuid.New()

	res := NewResolution()
	res.Default = ReadWrite
	res.SetCategory(open, "open", ReadWrite)
	res.SetCategory(closed, "closed", None)

	if res.ReadAll() {
		t.Error("ReadAll() = true with a category below read")
	}
	if slugs := res.AllowedReadSlugs(); len(slugs) != 1 || slugs[0] != "open" {
		t.Errorf("AllowedReadSlugs() = %v, want [open]", slugs)
	}

	res.SetCategory(closed, "closed", ReadOnly)
	if !res.ReadAll() || res.AllowedReadSlugs() != nil {
		t.Error("every category readable should mean no filter")
	}
}
