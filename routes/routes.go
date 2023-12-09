package routes

import (
	"SHforum_backend/controllers"
	"SHforum_backend/logger"
	"SHforum_backend/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	_ "SHforum_backend/docs"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(2*time.Second, 1))
	// 注册swagger路由
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	userRouter := v1.Group("/user")
	{
		userRouter.POST("/signup", controllers.SignUpHandler)
		userRouter.POST("/login", controllers.LoginHandler)
		//发送邮箱验证码
		userRouter.GET("/sendcode", controllers.SendCodeHandler)
		userRouter.Use(middlewares.JWTAuthMiddleware())
		{
			//关注用户
			userRouter.POST("/follow", controllers.FollowHandler)
			//取消关注
			userRouter.POST("/unfollow", controllers.UnFollowHandler)
		}
	}
	communityRouter := v1.Group("/community")
	{
		communityRouter.GET("/", controllers.CommunityHandler)
		communityRouter.GET("/:id", controllers.CommunityDetailHandler)
	}

	v1.GET("/post/:id", controllers.GetPostDetailHandler)
	v1.GET("/posts", controllers.GetPostListHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts2", controllers.GetPostListHandler2)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		v1.POST("/post", controllers.CreatePostHandler)
		v1.POST("/vote", controllers.PostVoteController)
		v1.GET("/search", controllers.SearchPostHandler)
	}
	pprof.Register(r)
	r.GET("/ping", func(c *gin.Context) {
		//如果是登录的用户，此时经过中间件，已经判断当前请求头中是否携带有效的JWT
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
