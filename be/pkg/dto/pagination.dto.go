package dto

type PaginationQueryDto struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

type LimitOffsetPaginationQueryDto struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type CursorTimePaginationQueryDto struct {
	LimitOffsetPaginationQueryDto
	CursorTime string `query:"cursor_time"`
}

type CursorPaginationQueryDto struct {
	Limit      int    `query:"limit"`
	CursorTime string `query:"cursor_time"`
}

type PaginationResponseDto struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
}

type LimitOffsetPaginationResponseDto struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type CursorTimePaginationResponseDto struct {
	LimitOffsetPaginationResponseDto
	NextCursorTime string `json:"next_cursor_time,omitempty"`
	HasMore        bool   `json:"has_more"`
}

type CursorPaginationResponseDto struct {
	Limit          int    `json:"limit"`
	NextCursorTime string `json:"next_cursor_time,omitempty"`
	HasMore        bool   `json:"has_more"`
}
