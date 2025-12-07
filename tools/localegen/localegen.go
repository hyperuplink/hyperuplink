package localegen

import (
	"fmt"
	"path"
)

func LocaleGen() (err error) {
	var ts []string

	if ts, err = walkDir("."); err != nil {
		return err
	}

	fmt.Printf("%v\n", ts)

	var locEN map[string]string
	if locEN, err = loadLocale(path.Join(".", "locales", "en.toml")); err != nil {
		return err
	}

	processLocales(path.Join(".", "locales"), locEN, ts)

	return nil
}
