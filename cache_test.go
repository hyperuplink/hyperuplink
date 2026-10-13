package main

import (
	"testing"

	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicperms "xn--gckvb8fzb.com/hyperuplink/logic/helpers/perms"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func TestSettingsAreServedFromTheCacheUntilWritten(t *testing.T) {
	h := apiTest(t)
	repo := gh.Repositories(h.rt).Setting

	system, err := settingRepo.GetByID[setting.System](h.ctx, repo, "system")
	if err != nil {
		t.Fatal(err)
	}
	name := system.JSONValue.Name
	t.Cleanup(func() {
		system.JSONValue.Name = name
		if err := settingRepo.Update(h.ctx, repo, system); err != nil {
			t.Errorf("restore the board name: %v", err)
		}
	})

	if _, err = database.From(h.rt).Exec(h.ctx,
		`UPDATE settings SET json_value = jsonb_set(json_value, '{name}', '"Renamed by SQL"') WHERE id = 'system'`,
	); err != nil {
		t.Fatal(err)
	}
	cached, err := settingRepo.GetByID[setting.System](h.ctx, repo, "system")
	if err != nil {
		t.Fatal(err)
	}
	if cached.JSONValue.Name != name {
		t.Errorf("name after a write behind the repository = %q, want the cached %q", cached.JSONValue.Name, name)
	}

	cached.JSONValue.Name = "Renamed by the repository"
	if err = settingRepo.Update(h.ctx, repo, cached); err != nil {
		t.Fatal(err)
	}
	fresh, err := settingRepo.GetByID[setting.System](h.ctx, repo, "system")
	if err != nil {
		t.Fatal(err)
	}
	if fresh.JSONValue.Name != "Renamed by the repository" {
		t.Errorf("name after Update = %q", fresh.JSONValue.Name)
	}

	fresh.JSONValue.Name = "Mutated by a caller"
	again, err := settingRepo.GetByID[setting.System](h.ctx, repo, "system")
	if err != nil {
		t.Fatal(err)
	}
	if again.JSONValue.Name != "Renamed by the repository" {
		t.Errorf("a caller's change reached the cache: %q", again.JSONValue.Name)
	}
}

func TestCategoriesAndPermissionsAreServedFromTheCacheUntilWritten(t *testing.T) {
	h := apiTest(t)
	repos := gh.Repositories(h.rt)

	categoryNamed := func(slug string) category.Category {
		t.Helper()
		cats, err := repos.Category.All(h.ctx, common.QueryOptions{})
		if err != nil {
			t.Fatal(err)
		}
		for _, cat := range cats {
			if cat.Slug == slug {
				return cat
			}
		}
		t.Fatalf("no category %q", slug)
		return category.Category{}
	}

	general := categoryNamed("general")
	t.Cleanup(func() {
		if err := repos.Category.Update(h.ctx, &general); err != nil {
			t.Errorf("restore the category: %v", err)
		}
	})

	if _, err := database.From(h.rt).Exec(h.ctx,
		`UPDATE categories SET name = 'Renamed by SQL' WHERE slug = 'general'`,
	); err != nil {
		t.Fatal(err)
	}
	if got := categoryNamed("general").Name; got != general.Name {
		t.Errorf("category after a write behind the repository = %q, want the cached %q", got, general.Name)
	}
	if err := repos.Category.MoveDown(h.ctx, &general); err != nil {
		t.Fatal(err)
	}
	if got := categoryNamed("general").Name; got != "Renamed by SQL" {
		t.Errorf("category after a repository write = %q, want the row", got)
	}
	if err := repos.Category.MoveUp(h.ctx, &general); err != nil {
		t.Fatal(err)
	}

	level := func() byte {
		return logicperms.Resolve(h.ctx, h.rt, user.UserRole, nil).DefaultLevel()
	}
	if got := level(); got != permission.ReadWrite {
		t.Fatalf("default level = %d, want read-write", got)
	}
	if _, err := database.From(h.rt).Exec(h.ctx,
		`UPDATE permissions SET bits = B'100' WHERE group_id IS NULL AND category_id IS NULL AND deleted_at IS NULL`,
	); err != nil {
		t.Fatal(err)
	}
	if got := level(); got != permission.ReadWrite {
		t.Errorf("default level after a write behind the repository = %d, want the cached read-write", got)
	}
	h.setDefaultPermission(t, permission.LevelRead)
	if got := level(); got != permission.ReadOnly {
		t.Errorf("default level after Apply = %d, want read", got)
	}
}

func TestDeletingAGroupDropsItsGrantsFromTheCache(t *testing.T) {
	h := apiTest(t)

	if err := logicpermissions.GroupCreate(h.ctx, h.rt, &logicpermissions.GroupCreateInput{
		ID:   "temporary",
		Name: "Temporary",
	}); err != nil {
		t.Fatal(err)
	}
	if err := logicpermissions.Apply(h.ctx, h.rt, &logicpermissions.ApplyInput{
		GroupID:    "temporary",
		CategoryID: h.categories["general"].String(),
		Level:      permission.LevelReadWriteModerate,
	}); err != nil {
		t.Fatal(err)
	}

	member := []string{"temporary"}
	if got := logicperms.Resolve(h.ctx, h.rt, user.UserRole, member).LevelBySlug("general"); got != permission.Moderate {
		t.Fatalf("level with the group = %d, want moderate", got)
	}

	if err := logicpermissions.GroupDestroy(h.ctx, h.rt, &logicpermissions.GroupDestroyInput{ID: "temporary"}); err != nil {
		t.Fatal(err)
	}
	res := logicperms.Resolve(h.ctx, h.rt, user.UserRole, member)
	if got := res.LevelBySlug("general"); got != res.DefaultLevel() {
		t.Errorf("level after the group is gone = %d, want the default %d", got, res.DefaultLevel())
	}
}
