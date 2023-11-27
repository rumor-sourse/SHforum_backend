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
	// 注册业务路由
	v1.POST("/signup", controllers.SignUpHandler)
	//登录路由
	v1.POST("/login", controllers.LoginHandler)

	v1.GET("/community", controllers.CommunityHandler)
	v1.GET("/community/:id", controllers.CommunityDetailHandler)
	v1.GET("/post/:id", controllers.GetPostDetailHandler)
	v1.GET("/posts", controllers.GetPostListHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts2", controllers.GetPostListHandler2)
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件

	{
		v1.POST("/post", controllers.CreatePostHandler)
		v1.POST("/vote", controllers.PostVoteController)
	}
	pprof.Register(r)
	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录的用户，此时经过中间件，已经判断当前请求头中是否携带有效的JWT
		c.String(http.StatusOK, "pong")
	})
	r.GET("/hello", func(c *gin.Context) {
		//如果是登录的用户，此时经过中间件，已经判断当前请求头中是否携带有效的JWT
		c.String(http.StatusOK, "world")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
