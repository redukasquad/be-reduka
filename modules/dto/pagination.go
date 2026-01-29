package dto

// ListQueryParams - untuk parsing query string ?q=&page=&perPage=
type ListQueryParams struct {
    Q       string `form:"q"`       // Search keyword
    Page    int    `form:"page"`    // Default: 1
    PerPage int    `form:"perPage"` // Default: 10
}
// SetDefaults - set default values jika tidak ada
func (p *ListQueryParams) SetDefaults() {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.PerPage < 1 {
        p.PerPage = 10
    }
    if p.PerPage > 100 {
        p.PerPage = 100  // Max limit untuk prevent abuse
    }
}
// GetOffset - hitung offset untuk SQL LIMIT
func (p *ListQueryParams) GetOffset() int {
    return (p.Page - 1) * p.PerPage
}
// PaginatedResponse - generic response dengan pagination info
type PaginatedResponse[T any] struct {
    Data       []T   `json:"data"`
    Page       int   `json:"page"`
    PerPage    int   `json:"perPage"`
    TotalItems int64 `json:"totalItems"`
    TotalPages int   `json:"totalPages"`
}
// NewPaginatedResponse - helper untuk membuat paginated response
func NewPaginatedResponse[T any](data []T, page, perPage int, totalItems int64) PaginatedResponse[T] {
    totalPages := int(totalItems) / perPage
    if int(totalItems)%perPage > 0 {
        totalPages++
    }
    
    return PaginatedResponse[T]{
        Data:       data,
        Page:       page,
        PerPage:    perPage,
        TotalItems: totalItems,
        TotalPages: totalPages,
    }
}