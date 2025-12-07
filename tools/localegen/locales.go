package localegen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func loadLocale(path string) (loc map[string]string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, &loc); err != nil {
		return nil, err
	}

	return loc, nil
}

func saveLocale(path string, loc map[string]string) (err error) {
	data, err := toml.Marshal(loc)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	// fmt.Printf("Would write to file %s:\n%s\n\n", path, data)

	return nil
}

func processLocales(
	dir string,
	ref map[string]string,
	ts []string,
) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".toml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		var loc, newloc map[string]string
		loc, err = loadLocale(path)

		for _, t := range ts {
			if val, ok := loc[t]; !ok ||
				val == "" ||
				(strings.HasPrefix(val, "#{{") && strings.HasSuffix(val, "}}#")) ||
				(strings.HasPrefix(val, "!{{") && strings.HasSuffix(val, "}}!")) {
				var setval string = "!{{" + t + "}}!"
				if refval, refok := ref[t]; refok && refval != "" {
					setval = "#{{" + refval + "}}#"
				}
				newloc[t] = setval
			}
		}

		if err = saveLocale(path, newloc); err != nil {
			return err
		}
	}

	return nil
}
