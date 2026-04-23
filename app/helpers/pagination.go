package helpers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Metadata struct {
	PageSize   int   `json:"page_size"`
	NextCursor int64 `json:"next_cursor,omitempty"`
	HasNext    bool  `json:"has_next"`
}

type Pagination struct {
	Page     *int
	PageSize *int
	Offset   int
	LastID   int64
	Metadata
}

func (p *Pagination) GetParams(c *gin.Context, v *Validator) {
	p.PageSize = GetIntParam(c, v, "pageSize")

	// Lấy last_id cho Keyset Pagination (Cursor)
	lastIDStr := c.Query("last_id")
	if lastIDStr != "" {
		id, err := strconv.ParseInt(lastIDStr, 10, 64)
		if err == nil {
			p.LastID = id
		}
	}

	if p.PageSize == nil || *p.PageSize <= 0 {
		defaultPageSize := 10
		p.PageSize = &defaultPageSize
	}
}

// CalculateMetadata cho kiểu Cursor-based (Limit + 1)
func (p *Pagination) CalculateMetadata(hasNext bool, nextCursor int64) {
	p.Metadata = Metadata{
		PageSize:   *p.PageSize,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}
}
