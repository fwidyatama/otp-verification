package route

import (
	"github.com/fwidyatama/e-recruitment/app/controller"
	"github.com/gin-gonic/gin"
)

func ServeRoute() {
	router := gin.Default()

	users := router.Group("/api/users")
	{
		users.POST("/", controller.RegisterUser)
		users.GET("/", controller.GetAllUser)
		users.POST("/verify", controller.VerifyOtp)
	}
	_ = router.Run()
}
