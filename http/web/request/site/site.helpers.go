package site

import "strings"

func (s *Site) Quote(text string) (quote string) {
	if text == "" {
		return text
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = "> " + line
	}
	return (strings.Join(lines, "\n") + "\n")
}
