package pagination

import "math"

type PageItemType int

const (
	RegularPage PageItemType = iota
	FirstPage                // i.e. left arrow
	LastPage                 // i.e. right arrow
)

type Page struct {
	Number   int
	Offset   int
	Size     int
	Active   bool
	ItemType PageItemType
}

func ComputePages(maxPages int, offset int, pageSize int, totalSize int) []Page {
	var pages []Page

	pagesTotal := int(math.Ceil(float64(totalSize) / float64(pageSize)))
	pagesRest := int(math.Ceil(float64(totalSize-offset) / float64(pageSize)))
	pagesBefore := pagesTotal - pagesRest

	pagesToLeft := maxPages / 2
	pagesToRight := maxPages - pagesToLeft // pagesToRight also includes the actual page

	if pagesToLeft > pagesBefore {
		pagesToRight += pagesToLeft - pagesBefore
		pagesToLeft = pagesBefore
	}

	if pagesToRight > pagesRest {
		pagesToLeft += pagesToRight - pagesRest
		pagesToRight = pagesRest
	}

	for i := range pagesToLeft {
		pageNumber := (pagesTotal - pagesRest) - i
		if pageNumber <= 0 {
			break
		}
		offset := (pageNumber - 1) * pageSize
		pages = append([]Page{{
			Number:   pageNumber,
			Offset:   offset,
			Size:     pageSize,
			Active:   false,
			ItemType: RegularPage,
		}}, pages...)
	}

	for i := range pagesToRight {
		pageNumber := (pagesTotal - pagesRest) + i + 1
		offset := (pageNumber - 1) * pageSize
		size := pageSize
		if i == pagesRest-1 {
			size = totalSize - offset
		}
		active := false
		if i == 0 {
			active = true
		}
		pages = append(pages, Page{
			Number:   pageNumber,
			Offset:   offset,
			Size:     size,
			Active:   active,
			ItemType: RegularPage,
		})
	}

	return pages
}

func GetActualPage(maxPages int, offset int, pageSize int, totalSize int) Page {
	pagesTotal := int(math.Ceil(float64(totalSize) / float64(pageSize)))
	pagesRest := int(math.Ceil(float64(totalSize-offset) / float64(pageSize)))
	pageNumber := (pagesTotal - pagesRest) + 1
	size := pageSize
	if pageNumber == pagesTotal {
		size = totalSize - offset
	}
	return Page{
		Number:   pageNumber,
		Offset:   offset,
		Size:     size,
		Active:   true,
		ItemType: RegularPage,
	}
}

func IncludeArrowPages(pages []Page, maxPages int, offset int, pageSize int, totalSize int) []Page {
	pagesTotal := int(math.Ceil(float64(totalSize) / float64(pageSize)))

	size := pageSize
	if pagesTotal == 1 {
		size = totalSize - offset
	}

	firstPage := Page{
		Number:   1,
		Offset:   0,
		Size:     size,
		Active:   false,
		ItemType: FirstPage,
	}

	lastPage := Page{
		Number:   pagesTotal,
		Offset:   (pagesTotal - 1) * pageSize,
		Size:     totalSize - offset,
		Active:   false,
		ItemType: LastPage,
	}

	result := append([]Page{firstPage}, pages...)
	result = append(result, lastPage)
	return result
}
