package response

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
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id"`
	Status      int32  `json:"status"`
}

type PostDetailResponse struct {
	AuthorName         string `json:"author_name"`
	VoteNum            int64  `json:"vote_num"`
	*PostResponse      `json:"post"`
	*CommunityResponse `json:"community"`
}
