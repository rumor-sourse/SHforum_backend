package controllers

import (
	"SHforum_backend/internal/logic"
	"SHforum_backend/internal/models"
	"SHforum_backend/internal/settings"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/url"
	"strconv"
)

type PostController struct {
}

// CreatePostHandler 创建帖子
func (pc *PostController) CreatePostHandler(c *gin.Context) {
	// 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// CanEditPostHandler 判断当前用户是否可以编辑帖子
func (pc *PostController) CanEditPostHandler(c *gin.Context) {
	//  获取参数（帖子id）
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id获取帖子数据
	data, err := logic.GetPostById(postID)
	if err != nil {
		zap.L().Error("GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 从请求中获取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 判断是否有权限
	if data.AuthorID != userID {
		ResponseError(c, CodeNotPermission)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// UpdatePostHandler 更新帖子
func (pc *PostController) UpdatePostHandler(c *gin.Context) {
	// 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdatePost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 更新帖子
	if err := logic.UpdatePost(p); err != nil {
		zap.L().Error("UpdatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// DeletePostHandler 删除帖子
func (pc *PostController) DeletePostHandler(c *gin.Context) {
	// 获取参数及参数校验
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("DeletePost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 删除帖子
	if err := logic.DeletePost(postID); err != nil {
		zap.L().Error("DeletePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 查询某个贴子详情
func (pc *PostController) GetPostDetailHandler(c *gin.Context) {
	// 获取参数（帖子id）
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id获取帖子数据
	data, err := logic.GetPostById(postID)
	if err != nil {
		zap.L().Error("GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 帖子列表接口
// @Summary 帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [get]
func (pc *PostController) GetPostListHandler(c *gin.Context) {
	//GET请求参数： /api/v1/posts?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取帖子列表数据
	data, err := logic.GetPostList(p)
	if err != nil {
		zap.L().Error("GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// SearchPostHandler 搜索帖子
func (pc *PostController) SearchPostHandler(c *gin.Context) {
	// 获取参数
	keyWord := c.Query("keyword")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 搜索帖子
	data, err := logic.SearchPost(keyWord, page)
	if err != nil {
		zap.L().Error("SearchPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// PostVoteHandler 为帖子投票
func (pc *PostController) PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	// 获取当前请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 业务处理
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// GetCommentByPostIdHandler 根据帖子id获取评论列表
func (pc *PostController) GetCommentByPostIdHandler(c *gin.Context) {
	pid := c.Param("id")
	postID, err := strconv.ParseInt(pid, 10, 64)
	//查询数据
	data, err := logic.GetCommentList(postID)
	if err != nil {
		zap.L().Error("GetCommentList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSuccess(c, data)
}

// GetHotCommentByPostIdHandler 根据帖子id获取热门评论
func (pc *PostController) GetHotCommentByPostIdHandler(c *gin.Context) {
	//获取参数
	pid := c.Param("id")
	postID, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		zap.L().Error("GetHotCommentByPostIdHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//查询数据
	data, err := logic.GetHotComment(postID)
	if err != nil {
		zap.L().Error("GetHotCommentList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回数据
	ResponseSuccess(c, data)
}

func (pc *PostController) SharePostHandler(c *gin.Context) {
	longUrl := GetBaseShareUrl(c)
	shortUrl, err := logic.Share(longUrl)
	if err != nil {
		zap.L().Error("GetShortURL failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, shortUrl)
}

func (pc *PostController) ShareCommentHandler(c *gin.Context) {
	longUrl := GetBaseShareUrl(c)
	shortUrl, err := logic.Share(longUrl)
	if err != nil {
		zap.L().Error("GetShortURL failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, shortUrl)
}

func GetBaseShareUrl(c *gin.Context) string {
	//获取当前url
	currentURL := c.Request.URL
	longUrl := &url.URL{
		Scheme: "http",
		Host:   settings.Conf.Host + ":" + strconv.FormatInt(int64(settings.Conf.Port), 10),
		Path:   currentURL.Path[:len(currentURL.Path)-len("/share")],
	}
	return longUrl.String()
}
