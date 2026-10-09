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

func TestAPIRoutesProduceTheVersionedPrefix(t *testing.T) {
	Use()

	if got, want := route.For("APIV1").AsURL(), "api/v1"; got != want {
		t.Errorf("APIV1.AsURL() = %q, want %q", got, want)
	}

	got := route.For("APIV1Topic").Fill(map[string]string{
		"topic": "0192b1c2-7d3e-7f40-8a51-00000000beef",
	}).AsURL()
	if want := "api/v1/topics/0192b1c2-7d3e-7f40-8a51-00000000beef"; got != want {
		t.Errorf("APIV1Topic.Fill().AsURL() = %q, want %q", got, want)
	}

	if !route.CollidesWithRoute("/api/anything") {
		t.Error("CollidesWithRoute(\"/api/anything\") = false, want true")
	}
}
