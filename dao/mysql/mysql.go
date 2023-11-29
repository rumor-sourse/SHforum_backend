package mysql

import (
	"SHforum_backend/models"
	"SHforum_backend/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init(cfg *settings.MySQLConfig) (err error) {
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
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		zap.L().Error("auto migrate user tables failed", zap.Error(err))
		return err
	}
	err = db.AutoMigrate(&models.Community{})
	if err != nil {
		zap.L().Error("auto migrate community tables failed", zap.Error(err))
		return err
	}
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		zap.L().Error("auto migrate post tables failed", zap.Error(err))
		return err
	}
	return
}
