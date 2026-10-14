package main

import (
	"path/filepath"
	"testing"

	"xn--gckvb8fzb.com/glides/services/database"
)

func TestCreateUserStoresTheMatchedLocale(t *testing.T) {
	h := apiTest(t)

	rt, err := newRuntime("file://" + filepath.Join(h.dir, "hyperuplink.toml"))
	if err != nil {
		t.Fatal(err)
	}
	rt.AddEmbed("migrations", &embedMigrations)
	database.From(rt).SetMigrations(rt.GetEmbed("migrations"))

	if err = createUser(rt, `{"username":"cliuser","email":"cliuser@example.org","password":"hyperhyper!","language":"de-AT"}`); err != nil {
		t.Fatalf("createUser: %v", err)
	}

	usr, err := h.reloadUser("cliuser")
	if err != nil {
		t.Fatalf("the created user can't be loaded: %v", err)
	}
	if usr.Language != "de" {
		t.Errorf("language = %q, want de", usr.Language)
	}
	if !usr.ConfirmedAt.Valid {
		t.Error("the created user isn't activated")
	}
}
