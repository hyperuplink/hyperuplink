package helpers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func GetPaths(c fiber.Ctx) (
	pathName string,
	absPath string,
	relRoot string,
) {
	cR := c.Route()
	absPath = cR.Path
	splitAbsPath := strings.Split(absPath, "/")
	lenSplitAbsPath := len(splitAbsPath)
	if lenSplitAbsPath > 0 {
		pathName = splitAbsPath[(lenSplitAbsPath - 1)]
	} else {
		pathName = ""
	}

	reqPath := c.Path()
	depth := 0
	if trimmed := strings.Trim(reqPath, "/"); trimmed != "" {
		depth = strings.Count(trimmed, "/") + 1
	}
	if depth > 0 && !strings.HasSuffix(reqPath, "/") {
		depth--
	}

	relRoot = "./"
	for i := 0; i < depth; i++ {
		relRoot += "../"
	}

	return pathName, absPath, relRoot
}
