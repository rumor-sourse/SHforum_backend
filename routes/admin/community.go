package admin

import (
	"SHforum_backend/interval/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterCommunityRouter(api *gin.RouterGroup) {
	communityController := new(controllers.CommunityController)
	communityRouter := api.Group("/community")
	{
		communityRouter.GET("/", communityController.GetCommunitiesHandler)
		communityRouter.GET("/:id", communityController.GetCommunitiesDetailHandler)
	}
}
