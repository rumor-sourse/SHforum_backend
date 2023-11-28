package mysql

import (
	"SHforum_backend/models"
	"gorm.io/gorm/clause"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	//insert into post(post_id, title, content, author_id, community_id) values(?,?,?,?,?)
	db.Debug().Create(p)
	return
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(pid int64) (data *models.Post, err error) {
	//select post_id, title, content, author_id, community_id, create_time from post where post_id = ?
	data = new(models.Post)
	db.Debug().Select("post_id", "title", "content", "author_id", "community_id", "create_time").Where("post_id = ?", pid).First(data)
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.Post, err error) {
	//select post_id, title, content, author_id, community_id, create_time from post order by create_time desc limit ?,?
	data = make([]*models.Post, 0, 2)
	result := db.Debug().Select("post_id", "title", "content", "author_id", "community_id", "create_time").Order("create_time desc").Limit(int(size)).Offset(int((page - 1) * size)).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return
}

// GetPostListByIDs 根据多个ID查询帖子列表数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	/*	sqlStr := `select post_id, title, content, author_id, community_id, create_time from post where post_id in (?) order by FIND_IN_SET(post_id, ?)`
		query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
		if err != nil {
			return nil, err
		}
		query = db.Rebind(query)
		err = db.Select(&postList, query, args...)*/
	postList = make([]*models.Post, len(ids))
	db.Debug().Select("post_id", "title", "content", "author_id", "community_id", "create_time").
		Where("post_id IN ?", ids).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(post_id,?)", Vars: []interface{}{ids}, WithoutParentheses: true},
	}).Find(&postList)
	return
}
