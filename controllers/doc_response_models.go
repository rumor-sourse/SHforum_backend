package controllers

import (
	"SHforum_backend/models/response"
)

//专门用来放接口文档用到的odel

type _ResponsePostList struct {
	Code    ResCode                  `json:"code"`
	Message string                   `json:"message"`
	Data    []*response.PostResponse `json:"data"`
}
