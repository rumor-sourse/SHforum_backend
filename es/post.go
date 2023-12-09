package es

import (
	"SHforum_backend/models"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"strconv"
)

func CreatePostIndex(post models.Post) error {
	_, err := client.Index().
		Index("post").                  //设置索引名称
		Id(strconv.Itoa(int(post.ID))). //设置文档id
		BodyJson(post).                 //设置数据
		Do(context.Background())        //执行请求
	zap.L().Info("es.CreatePostIndex(post) success")
	if err != nil {
		return err
	}
	return nil
}

func SearchPostIndex(keyword string, page int) ([]models.Post, error) {
	size := 15
	res, err := client.Search("post").
		Query(elastic.NewMultiMatchQuery(keyword, "title", "content").Analyzer("whitespace")).
		From((page - 1) * size).
		Size(size).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	// 解析结果
	var posts []models.Post
	for _, hit := range res.Hits.Hits {
		var post models.Post
		err := json.Unmarshal(hit.Source, &post)
		if err != nil {
			zap.L().Error("json.Unmarshal(hit.Source, &post) failed",
				zap.Error(err))
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}
