package common

import "math"

type PaginationRequest struct {
	PerPage int `json:"per_page" form:"per_page"`
	Page    int `json:"page" form:"page"`
}

func NewPaginationRequet(page int, perPage int) PaginationRequest {
	paginationRequest := PaginationRequest{
		PerPage: perPage,
		Page:    page,
	}
	if perPage == 0 {
		paginationRequest.PerPage = 20
	} else {
		paginationRequest.PerPage = perPage
	}
	if page == 0 {
		paginationRequest.Page = 1
	} else {
		paginationRequest.Page = page
	}

	return paginationRequest
}

type Pagination struct {
	PaginationRequest
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

func NewPagination(page, perPage, totalItems int) *Pagination {
	return &Pagination{
		PaginationRequest: PaginationRequest{
			PerPage: perPage,
			Page:    page,
		},
		TotalPages: int(math.Ceil(float64(totalItems) / float64(perPage))),
		TotalItems: totalItems,
	}
}
