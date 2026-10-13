package helpers

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/intnat"
)

const DefaultLocale = "en"

var Locales = []string{"de", "en", "es", "fr", "it", "ro"}

var LocaleNames = map[string]string{
	"de": "Deutsch",
	"en": "English",
	"es": "Español",
	"fr": "Français",
	"it": "Italiano",
	"ro": "Română",
}

func Locale(rt *runtime.Runtime, accepts ...string) string {
	return intnat.From(rt).Bundle.MatchAvailableLocale(accepts...)
}
