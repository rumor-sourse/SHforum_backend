package admin

import (
	"SHforum_backend/interval/controllers"
	"SHforum_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(api *gin.RouterGroup) {
	userController := new(controllers.UserController)
	userRouter := api.Group("/user")
	{
		userRouter.POST("/signup", userController.SignUpHandler)
		userRouter.POST("/login", userController.LoginHandler)
		// 发送邮箱验证码
		userRouter.GET("/sendcode", userController.SendCodeHandler)
		userRouter.Use(middlewares.JWTAuthMiddleware())
		{
			// 关注用户
			userRouter.POST("/follow", userController.FollowHandler)
			// 取消关注
			userRouter.POST("/unfollow", userController.UnFollowHandler)
			// 更新用户信息
			userRouter.POST("/update/:id", userController.UpdateUserByIDHandler)
			// 私信
			userRouter.POST("/message/:id", userController.SendMessageHandler)
		}
	}
}
