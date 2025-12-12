package route

import "testing"

func TestRouteFill(t *testing.T) {
	r := For("CategoriesForumsTopics")
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
