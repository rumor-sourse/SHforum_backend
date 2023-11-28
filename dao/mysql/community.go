package mysql

import (
	"SHforum_backend/models"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (data []models.Community, err error) {
	//select community_id, community_name from community
	result := db.Debug().Select("community_id", "community_name").Find(&data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Error("community list not found", zap.Error(result.Error))
		return nil, result.Error
	}
	return
}

// GetCommunityDetailByID 根据社区ID获取社区详情
func GetCommunityDetailByID(id int64) (data *models.CommunityDetail, err error) {
	data = new(models.CommunityDetail)
	//select community_id, community_name, introduction, create_time from community where community_id = ?
	result := db.Debug().Select("community_id", "community_name", "introduction", "create_time").Where("community_id = ?", id).First(data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Warn("there is no community in table community")
		return nil, result.Error
	}
	return
}
