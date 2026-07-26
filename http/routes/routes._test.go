package routes

import (
	"testing"

	"xn--gckvb8fzb.com/glides/http/route"
)

func TestRouteFill(t *testing.T) {
	Use()

	r := route.For("CategoriesForumsTopics")
	fr := r.Fill(
		map[string]string{
			"categories": "hosting",
			"forums":     "aws",
			"topics":     "evil",
		},
	)

	if fr.AsURL() != "_hosting/aws/evil" {
		t.Errorf("Fill() incorrect: %s", fr.AsURL())
	}
}

func TestUseSetsTheDefaultEnvironmentTitle(t *testing.T) {
	Use()

	if got, want := route.NewEnv().Title, DEFAULT_TITLE; got != want {
		t.Errorf("NewEnv().Title = %q, want %q", got, want)
	}
}

func TestEveryRouteResolves(t *testing.T) {
	Use()

	for id := range Routes {
		if route.For(id).Len() == 0 {
			t.Errorf("route %q does not resolve through route.For", id)
		}
	}
}
