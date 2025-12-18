package helpers

import (
	"math"
	"strings"

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
