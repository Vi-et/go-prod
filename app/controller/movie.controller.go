package controller

import (
	"net/http"
	"strconv"

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
		return
	}

	// 5. Trả về kết quả JSON kèm metadata
	c.JSON(http.StatusOK, gin.H{
		"metadata": metadata,
		"movies":   movies,
	})
}

func (mc *MovieController) GetController(c *gin.Context) {
	id := c.Param("id")

	movie, err := model.Movie{}.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

func (mc *MovieController) CreateController(c *gin.Context) {
	var movie model.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := movie.Create(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}

func (mc *MovieController) UpdateController(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var inputUpdate model.InputUpdateMovie

	if err := c.ShouldBindJSON(&inputUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := model.Movie{}.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = movie.Update(&inputUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}
