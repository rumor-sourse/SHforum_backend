package admin

import (
	"SHforum_backend/interval/controllers"
	"SHforum_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterPostRouter(api *gin.RouterGroup) {
	postController := new(controllers.PostController)
	//根据时间或分数获取帖子列表
	api.GET("/posts", postController.GetPostListHandler)
	api.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件
	{
		//根据内容搜索博客，通过es实现
		api.GET("/search", postController.SearchPostHandler)
	}
	postRouter := api.Group("/post", middlewares.JWTAuthMiddleware())
	{
		//查询某个贴子详情
		postRouter.GET("/:id", postController.GetPostDetailHandler)
		//创建贴子
		postRouter.POST("/create", postController.CreatePostHandler)
		//判断某个贴子是否可以编辑
		postRouter.GET("/canEdit/:id", postController.CanEditPostHandler)
		//更新贴子
		postRouter.POST("/update/:id", postController.UpdatePostHandler)
		//删除贴子
		postRouter.POST("/delete/:id", postController.DeletePostHandler)
		//为某个贴子投票
		postRouter.POST("/vote", postController.PostVoteHandler)
		//获取某个贴子的评论列表
		postRouter.GET("/comments/:id", postController.GetCommentByPostIdHandler)
		//获取某个贴子的热评
		postRouter.GET("/hotcomment/:id", postController.GetHotCommentByPostIdHandler)
	}
}
