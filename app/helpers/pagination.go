package helpers

import (
	"math"

	"github.com/gin-gonic/gin"
)

type Metadata struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	PageSize     int   `json:"page_size,omitempty"`
	FirstPage    int   `json:"first_page,omitempty"`
	LastPage     int   `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

type Pagination struct {
	Page     *int
	PageSize *int
	Offset   int
	Metadata
}

func (p *Pagination) GetParams(c *gin.Context, v *Validator) {
	p.Page = GetIntParam(c, v, "page")
	p.PageSize = GetIntParam(c, v, "pageSize")

	if p.Page == nil || *p.Page <= 0 {
		defaultPage := 1
		p.Page = &defaultPage
	}
	if p.PageSize == nil || *p.PageSize <= 0 {
		defaultPageSize := 10
		p.PageSize = &defaultPageSize
	}
	p.Offset = (*p.Page - 1) * *p.PageSize
}

// CalculateMetadata giúp tính toán các thông số phân trang nhanh chóng
func (p *Pagination) CalculateMetadata(totalRecords int64) {
	if totalRecords == 0 {
		return
	}
	p.Metadata = Metadata{
		CurrentPage:  *p.Page,
		PageSize:     *p.PageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(*p.PageSize))),
		TotalRecords: totalRecords,
	}
}
