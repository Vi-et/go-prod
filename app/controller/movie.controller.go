package controller

import (
	"net/http"
	"strconv"

	"go-production/app/helpers"
	"go-production/global"
	"go-production/internal/model"

	"github.com/gin-gonic/gin"
)

type MovieController struct{}

func NewMovieController() *MovieController {
	return &MovieController{}
}

// ListMovies xử lý GET /v1/movies?page=1&page_size=10
func (mc *MovieController) ListMovies(c *gin.Context) {
	// 1. Lấy thông tin phân trang từ query string
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var movies []model.Movie
	var totalRecords int64

	// 2. Đếm tổng số bản ghi
	if err := global.DB.Model(&model.Movie{}).Count(&totalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không thể đếm số lượng phim"})
		return
	}

	// 3. Truy vấn dữ liệu với phân trang
	offset := (page - 1) * pageSize
	if err := global.DB.Limit(pageSize).Offset(offset).Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không thể lấy danh sách phim"})
		return
	}

	// 4. Tính toán Metadata bằng hàm global
	metadata := helpers.CalculateMetadata(totalRecords, page, pageSize)

	// 5. Trả về kết quả JSON kèm metadata
	c.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"movies":   movies,
	})
}
