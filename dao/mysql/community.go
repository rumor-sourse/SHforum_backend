package mysql

import (
	"SHforum_backend/models"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (data []*models.Community, err error) {
	//select id, name from community
	result := db.Debug().Select("id", "name").Find(&data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Error("community list not found", zap.Error(result.Error))
		return nil, result.Error
	}
	return
}

// GetCommunityDetailByID 根据社区ID获取社区详情
func GetCommunityDetailByID(id int64) (community *models.Community, err error) {
	//select community_id, community_name, introduction from community where community_id = ?
	community = new(models.Community)
	result := db.Debug().Select("id", "name", "introduction").Where("id = ?", id).First(&community)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Warn("there is no community in table community")
		return nil, result.Error
	}
	return
}
