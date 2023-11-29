package response

import "SHforum_backend/models"

type UserResponse struct {
	UserID int64  `json:"id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

type CommunityResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CommunityDetailResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

type PostResponse struct {
	AuthorName        string `json:"author_name"`
	VoteNum           int64  `json:"vote_num"`
	*models.Post      `json:"post"`
	*models.Community `json:"community"`
}
