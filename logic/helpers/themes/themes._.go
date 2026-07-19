package themes

import (
	"io/fs"
	"sort"

	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func getPackages(
	rt *runtime.Runtime,
	dir string,
	file string,
) (names []string, err error) {
	var entries []fs.DirEntry
	entries, err = fs.ReadDir(rt.Embeds["static"], dir)
	if err != nil {
		return nil, err
	}

	names = make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			if _, serr := fs.Stat(
				rt.Embeds["static"],
				dir+"/"+entry.Name()+"/"+file,
			); serr != nil {
				continue
			}
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	return names, nil
}

func GetThemes(
	rt *runtime.Runtime,
) (names []string, err error) {
	return getPackages(rt, "static/themes", "theme.css")
}

func GetColorschemes(
	rt *runtime.Runtime,
) (names []string, err error) {
	return getPackages(rt, "static/colorschemes", "colorscheme.css")
}
