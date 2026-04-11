package initialize

import (
	"go-production/app/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Khởi tạo controller
	movieCtrl := controller.NewMovieController()

	v1 := r.Group("/v1")
	{
		v1.GET("/movies", movieCtrl.ListMovies)
	}

	return r
}
