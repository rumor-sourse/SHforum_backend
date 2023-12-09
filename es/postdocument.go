package es

type PostDocument struct {
	ID          uint   `json:"post_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id"`
}
