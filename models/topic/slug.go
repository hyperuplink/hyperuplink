package topic

import (
	"regexp"
	"strings"
)

func (m *Topic) SetSlugFromName() {
	name := m.Name

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	re := regexp.MustCompile("[^a-z0-9-]")
	name = re.ReplaceAllString(name, "")

	m.Slug = name
}
