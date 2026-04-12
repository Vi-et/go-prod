package initialize

import (
	"expvar"
	"go-production/app/controller"
	"go-production/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Apply metrics middleware
	r.Use(middleware.MetricsMiddleware())

	// Expose expvar metrics
	r.GET("/debug/vars", gin.WrapH(expvar.Handler()))

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
