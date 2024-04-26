package routes

import (
	"SHforum_backend/controllers"
	_ "SHforum_backend/docs"
	"SHforum_backend/logger"
	"SHforum_backend/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func SetUp(mode string) *gin.Engine {
	//设置成发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true))
	middlewares.RateLimitMiddleware(time.Microsecond*time.Duration(200), 20000)
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
	//获取贴子列表
	//v1.GET("/posts", controllers.GetPostsHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts", controllers.GetPostListHandler)
	postRouter := v1.Group("/post", middlewares.JWTAuthMiddleware())
	{
		//查询某个贴子详情
		postRouter.GET("/:id", controllers.GetPostDetailHandler)
		//创建贴子
		postRouter.POST("/create", controllers.CreatePostHandler)
		//判断某个贴子是否可以编辑
		postRouter.GET("/canEdit/:id", controllers.CanEditPostHandler)
		//更新贴子
		postRouter.POST("/update/:id", controllers.UpdatePostHandler)
		//删除贴子
		postRouter.POST("/delete/:id", controllers.DeletePostHandler)
		//为某个贴子投票
		postRouter.POST("/vote", controllers.PostVoteController)
		//获取某个贴子的评论列表
		postRouter.GET("/comments", controllers.GetCommentByPostIdHandler)
		//获取某个贴子的热评
		postRouter.GET("/hotcomment/:id", controllers.GetHotCommentByPostIdHandler)
	}
	commentRouter := v1.Group("/comment", middlewares.JWTAuthMiddleware())
	{
		//创建评论
		commentRouter.POST("/add", controllers.CreateCommentHandler)
		//给某个评论点赞
		commentRouter.POST("/like", controllers.LikeCommentHandler)
		//判断某个评论是否可以编辑
		commentRouter.GET("/canEdit/:id", controllers.CanEditCommentHandler)
		//删除评论
		commentRouter.POST("/delete/:id", controllers.DeleteCommentHandler)
	}
	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		//根据内容搜索博客，通过es实现
		v1.GET("/search", controllers.SearchPostHandler)
	}
	pprof.Register(r)
	//测试
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	//404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
