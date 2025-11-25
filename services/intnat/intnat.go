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

func New() (*Intnat, error) {
	in := new(Intnat)

	in.Bundle = i18n.NewBundle(
		i18n.WithDefaultLocale("en"),
		i18n.WithLocales(
			"de",
			"en",
			"es",
			"fr",
			"it",
			"jp",
			"ro",
		),
		i18n.WithUnmarshaler(toml.Unmarshal),
	)

	return in, nil
}

func (in *Intnat) SetLocales(embedLocales *embed.FS) error {
	// var err error

	in.Locales = embedLocales

	return nil
}

func (in *Intnat) Startup() error {
	// var err error

	in.Bundle.LoadFS(in.Locales, "locales/*.toml")

	return nil
}

func (in *Intnat) Shutdown() error {
	// var err error

	return nil
}

func (in *Intnat) NewLocalizer(locale string) *i18n.Localizer {
	return in.Bundle.NewLocalizer(in.Bundle.MatchAvailableLocale(locale))
}
