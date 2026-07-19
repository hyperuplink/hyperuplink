package paging

import "math"

func Pages(total int64, perPage int) (pages int) {
	pages = int(math.Ceil(float64(total) / float64(perPage)))
	return pages
}
