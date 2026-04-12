package controller

import (
	"net/http"

	"go-production/app/filters"
	"go-production/app/helpers"
	"go-production/app/model"
	"go-production/global"

	"github.com/gin-gonic/gin"
)

type MovieController struct{}

func NewMovieController() *MovieController {
	return &MovieController{}
}

// ListMovies xử lý GET /v1/movies?page=1&page_size=10
func (mc *MovieController) ListMovies(c *gin.Context) {
	// 1. Lấy thông tin phân trang từ query string

	f := filters.MovieFilter{}
	v := helpers.Validator{}
	p := helpers.Pagination{}
	s := helpers.Ordering{}
	f.GetParams(c, &v)
	s.GetParams(c, &v)
	if !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": v.Errors})
		return
	}

	var movies []model.Movie
	var totalRecords int64

	// 2. Xây dựng câu lệnh query với filters
	query := global.DB.Model(&model.Movie{})
	query = f.Apply(query)

	// 3. Đếm tổng số bản ghi sau khi lọc
	if err := query.Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không thể thống kê số lượng phim"})
		return
	}

	// 4. Truy vấn dữ liệu với phân trang
	if err := query.Limit(p.PageSize).Offset(p.Offset).Order(s.Order).Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không thể lấy danh sách phim"})
		return
	}

	// 4. Tính toán Metadata bằng hàm global
	p.CalculateMetadata(totalRecords)

	// 5. Trả về kết quả JSON kèm metadata
	c.JSON(http.StatusOK, gin.H{
		"metadata": p.Metadata,
		"movies":   movies,
	})
}
