package initialize

import (
	"expvar"
	"go-production/app/controller"
	"go-production/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Apply rate limit middleware
	r.Use(middleware.RateLimitMiddleware())

	// Apply metrics middleware
	r.Use(middleware.MetricsMiddleware())

	// Apply authentication middleware globally
	r.Use(middleware.Authenticate())

	// Expose expvar metrics
	r.GET("/debug/vars", gin.WrapH(expvar.Handler()))

	// Khởi tạo controller
	movieCtrl := controller.NewMovieController()
	userCtrl := controller.NewUserController()

	v1 := r.Group("/v1")
	{
		// User routes
		v1.POST("/users", userCtrl.Register)
		v1.POST("/users/login", userCtrl.Login)
		v1.GET("/users/me", middleware.RequireAuthenticatedUser(), userCtrl.GetProfile)

		// Movie routes
		v1.GET("/movies", movieCtrl.ListController)
		v1.GET("/movies/:id", movieCtrl.GetController)
		v1.PATCH("/movies/:id", middleware.RequireAuthenticatedUser(), movieCtrl.UpdateController)
		v1.POST("/movies", middleware.RequireAuthenticatedUser(), movieCtrl.CreateController)
	}

	return r
}
