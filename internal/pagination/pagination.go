package pagination

import "math"

type Page struct {
	Number int
	Size   int
	Active bool
}

func ComputePages(maxPages int, offset int, pageSize int, totalSize int) []Page {
	var pages []Page

	pagesTotal := int(math.Ceil(float64(totalSize) / float64(pageSize)))
	pagesLeft := int(math.Ceil(float64(totalSize-offset) / float64(pageSize)))
	// lastPageSize :=

	for i := range pagesLeft { // but max maxPages
		active := false
		if i == 0 {
			active = true
		}
		pages = append(pages, Page{
			Number: (pagesTotal - pagesLeft) + i + 1,
			Size:   pageSize, // Except for the last one
			Active: active,
		})
	}

	return pages
}
