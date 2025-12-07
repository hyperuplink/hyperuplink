package localegen

import (
	"path"
)

func LocaleGen() (err error) {
	var ts []string

	if ts, err = walkDir("."); err != nil {
		return err
	}

	var locEN map[string]string
	if locEN, err = loadLocale(path.Join(".", "locales", "en.toml")); err != nil {
		return err
	}

	processLocales(path.Join(".", "locales"), locEN, ts)

	return nil
}
