package httpresponse

import "math"

type PaginationMeta struct {
	TotalDocs     int64 `json:"totalDocs"`
	Limit         int   `json:"limit"`
	Page          int   `json:"page"`
	TotalPages    int   `json:"totalPages"`
	PagingCounter int   `json:"pagingCounter"`
	HasPrevPage   bool  `json:"hasPrevPage"`
	HasNextPage   bool  `json:"hasNextPage"`
	PrevPage      *int  `json:"prevPage"`
	NextPage      *int  `json:"nextPage"`
}

type PaginatedData struct {
	Docs interface{} `json:"docs"`
	PaginationMeta
}

func BuildPaginationMeta(total int64, page, limit int) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	var prevPage, nextPage *int

	hasPrev := page > 1
	hasNext := page < totalPages

	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	if hasNext {
		n := page + 1
		nextPage = &n
	}

	return PaginationMeta{
		TotalDocs:     total,
		Limit:         limit,
		Page:          page,
		TotalPages:    totalPages,
		PagingCounter: (page-1)*limit + 1,
		HasPrevPage:   hasPrev,
		HasNextPage:   hasNext,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}
}