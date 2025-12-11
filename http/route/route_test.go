package route

import "testing"

func TestRouteFill(t *testing.T) {
	r := For("CategoriesForumsTopics")
	f := make(map[string]string)
	f["categories"] = "hosting"
	f["forums"] = "aws"
	f["topics"] = "evil"
	fr := r.Fill(f)

	if fr.AsURL() != "_hosting/aws/evil" {
		t.Errorf("Fill() incorrect: %s", fr.AsURL())
	}
}
