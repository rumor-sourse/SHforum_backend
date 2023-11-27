package mysql

import (
	"SHforum_backend/models"
	"SHforum_backend/settings"
	"testing"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "localhost",
		Port:         3306,
		User:         "root",
		Password:     "wszjdfs123456",
		DBName:       "shforum",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	err := CreatePost(&models.Post{
		ID:          10,
		AuthorID:    1,
		CommunityID: 1,
		Title:       "test",
		Content:     "test",
	})
	if err != nil {
		t.Errorf("CreatePost() failed, err:%v", err)
		return
	}
	t.Logf("CreatePost() success")
}
