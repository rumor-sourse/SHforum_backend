package mysql

import (
	"SHforum_backend/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(pid int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id = ?`
	if err := db.Get(data, sqlStr, pid); err != nil {
		return nil, err
	}
	return
}

func GetPostList(page, size int64) (data []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post order by create_time desc limit ?,?`
	data = make([]*models.Post, 0, 2)
	if err := db.Select(&data, sqlStr, (page-1)*size, size); err != nil {
		return nil, err
	}
	return
}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id in (?) order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
