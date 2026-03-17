package dto

type PaginationQueryDto struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type PaginationResponseDto struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
}
