package controller

import (
	"net/http"

	"go-production/app/filters"
	"go-production/app/helpers"
	"go-production/app/model"

	"github.com/gin-gonic/gin"
)

type MovieController struct{}

func NewMovieController() *MovieController {
	return &MovieController{}
}

// ListMovies xử lý GET /v1/movies?page=1&page_size=10
func (mc *MovieController) ListController(c *gin.Context) {
	// 1. Lấy thông tin phân trang từ query string

	f := filters.MovieFilter{}
	v := helpers.Validator{}
	o := helpers.Ordering{}
	p := helpers.Pagination{}
	f.GetParams(c, &v)
	o.GetParams(c, &v)
	p.GetParams(c, &v)
	if !v.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": v.Errors})
		return
	}

	metadata, movies, err := model.Movie{}.List(&f, &p, &o)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})

	}

	// 5. Trả về kết quả JSON kèm metadata
	c.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"movies":   movies,
	})
}
