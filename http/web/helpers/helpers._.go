package helpers

import (
	"io/fs"
	"math"
	"sort"
	"strings"

	"github.com/mrusme/hyperuplink/runtime"

	"github.com/gofiber/fiber/v3"
)

func GetPaths(c fiber.Ctx) (
	pathName string,
	absPath string,
	relRoot string,
) {
	origURL := c.OriginalURL()
	cR := c.Route()
	absPath = cR.Path
	splitAbsPath := strings.Split(absPath, "/")
	lenSplitAbsPath := len(splitAbsPath)
	if lenSplitAbsPath > 0 {
		pathName = splitAbsPath[(lenSplitAbsPath - 1)]
	} else {
		pathName = ""
	}
	parts := lenSplitAbsPath - 1

	if strings.HasSuffix(origURL, "/") {
		parts++
	}

	relRoot = "./"
	for i := 1; i < parts; i++ {
		relRoot += "../"
	}

	return pathName, absPath, relRoot
}

func GetNumberOfPages(total int64, perPage int) (pages int) {
	pages = int(math.Ceil(float64(total) / float64(perPage)))
	return pages
}

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
