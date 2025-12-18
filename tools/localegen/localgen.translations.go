package localegen

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/mrusme/hyperuplink/http/route"
)

var SUPPORTED_EXTS []string = []string{
	".go",
	".html",
	".eml",
}

func walkDir(dirpath string) (ts []string, err error) {
	err = filepath.Walk(dirpath, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			return e
		}

		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if slices.Index(SUPPORTED_EXTS, ext) == -1 {
			return nil
		}

		var extractedTs []string
		extractedTs, err = extractTranslations(path)
		if err != nil {
			return err
		}

		for _, xT := range extractedTs {
			if slices.Index(ts, xT) == -1 {
				ts = append(ts, xT)
			}
		}

		return nil
	})

	return ts, err
}

func extractTranslations(filepath string) (ts []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return ts, err
	}
	defer file.Close()

	re1 := regexp.MustCompile(`\.Ts{0,1}."(\w+)".`)
	re2 := regexp.MustCompile(`\.Form\.Input\s+"\w+"\s+"(\w+)"`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		submatches := re1.FindAllStringSubmatch(line, -1)
		for _, sm := range submatches {
			if len(sm) > 1 {
				ts = append(ts, sm[1])
			}
		}
		submatches = re2.FindAllStringSubmatch(line, -1)
		for _, sm := range submatches {
			if len(sm) > 1 {
				ts = append(ts, sm[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ts, err
	}

	for _, r := range route.Routes {
		title := r.AsTitle()
		if title != "" {
			ts = append(ts, title)
		}
	}

	return ts, nil
}
