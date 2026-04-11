package helpers

import (
	"go-production/globalStructs"
	"math"
)

// CalculateMetadata giúp tính toán các thông số phân trang nhanh chóng
func CalculateMetadata(totalRecords int64, page, pageSize int) globalStructs.Metadata {
	if totalRecords == 0 {
		return globalStructs.Metadata{}
	}
	return globalStructs.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
