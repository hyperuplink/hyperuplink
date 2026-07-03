package intnat

import (
	"embed"

	"github.com/kaptinlin/go-i18n"

	"github.com/pelletier/go-toml/v2"
)

type Intnat struct {
	Bundle  *i18n.I18n
	Locales *embed.FS
}

func New() (in *Intnat, err error) {
	in = new(Intnat)

	in.Bundle, err = i18n.NewBundle(
		i18n.WithDefaultLocale("en"),
		i18n.WithLocales(
			"de",
			"en",
			"es",
			"fr",
			"it",
			"ro",
		),
		i18n.WithUnmarshaler(toml.Unmarshal),
	)

	return in, err
}

func (in *Intnat) SetLocales(embedLocales *embed.FS) (err error) {
	in.Locales = embedLocales

	return nil
}

func (in *Intnat) Startup() (err error) {
	in.Bundle.LoadFS(in.Locales, "locales/*.toml")

	return nil
}

func (in *Intnat) Shutdown() (err error) {
	return nil
}

func (in *Intnat) NewLocalizer(locale string) *i18n.Localizer {
	return in.Bundle.NewLocalizer(in.Bundle.MatchAvailableLocale(locale))
}
