package es

import (
	"SHforum_backend/settings"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

var (
	client *elastic.Client
)

func InitEs(cfg *settings.EsConfig) (err error) {
	url := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	client, err = elastic.NewClient(
		//elastic 服务地址
		elastic.SetURL(url),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		return err
	}
	return nil
}
