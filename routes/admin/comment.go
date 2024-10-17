package admin

import (
	"SHforum_backend/internal/controllers"
	"SHforum_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterCommentController(api *gin.RouterGroup) {
	commentController := new(controllers.CommentController)
	commentRouter := api.Group("/comment", middlewares.JWTAuthMiddleware())
	{
		//创建评论
		commentRouter.POST("/add", commentController.CreateCommentHandler)
		//给某个评论点赞
		commentRouter.POST("/like", commentController.LikeCommentHandler)
		//判断某个评论是否可以编辑
		commentRouter.GET("/canEdit/:id", commentController.CanEditCommentHandler)
		//删除评论
		commentRouter.POST("/delete/:id", commentController.DeleteCommentHandler)
	}
}
