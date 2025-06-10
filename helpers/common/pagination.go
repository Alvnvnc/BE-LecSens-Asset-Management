package common

// QueryParams represents common query parameters for pagination and filtering
type QueryParams struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
	SortBy    string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	Search    string `json:"search" form:"search"`
}

// PaginationResponse represents the pagination metadata in API responses
type PaginationResponse struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// NewQueryParams creates a new QueryParams with default values
func NewQueryParams() *QueryParams {
	return &QueryParams{
		Page:      1,
		PageSize:  10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Validate validates and sets default values for QueryParams
func (qp *QueryParams) Validate() {
	if qp.Page <= 0 {
		qp.Page = 1
	}
	if qp.PageSize <= 0 {
		qp.PageSize = 10
	}
	if qp.PageSize > 100 {
		qp.PageSize = 100
	}
	if qp.SortBy == "" {
		qp.SortBy = "created_at"
	}
	if qp.SortOrder == "" {
		qp.SortOrder = "desc"
	}
}

// GetOffset calculates the offset for database queries
func (qp *QueryParams) GetOffset() int {
	return (qp.Page - 1) * qp.PageSize
}

// NewPaginationResponse creates a new PaginationResponse
func NewPaginationResponse(page, pageSize int, total int64) *PaginationResponse {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationResponse{
		Page:        page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}
