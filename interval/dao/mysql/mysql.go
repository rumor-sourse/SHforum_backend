package mysql

import (
	"SHforum_backend/interval/models"
	"SHforum_backend/interval/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Init(cfg *settings.MySQLConfig) (err error) {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		// 也可以使用MustConnect连接不成功就panic
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: dsn, // DSN data source name
		}), &gorm.Config{})
		if err != nil {
			zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
			return
		}
		sqlDB, err := db.DB()
		if err != nil {
			zap.L().Error("get DB failed", zap.Error(err))
		}
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		err = db.AutoMigrate(&models.User{},
			&models.Community{},
			&models.Post{},
			&models.Follow{},
			&models.Fan{},
			&models.Message{},
			&models.Comment{})
		if err != nil {
			zap.L().Error("auto migrate tables failed", zap.Error(err))
			return
		}
	})
	return
}
