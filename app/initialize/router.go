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
		v1.GET("/movies", movieCtrl.ListController)
		v1.GET("/movies/:id", movieCtrl.GetController)
		v1.PATCH("/movies/:id", movieCtrl.UpdateController)
		v1.POST("/movies", movieCtrl.CreateController)
	}

	return r
}
