package mysql

import (
	"SHforum_backend/models"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

func GetCommunityList() (data []models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	if err := db.Select(&data, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("there is no community in table community")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (data *models.CommunityDetail, err error) {
	data = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	if err := db.Get(data, sqlStr, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zap.L().Warn("there is no community in table community")
			err = ErrorInvalidID
		}
	}
	return data, err
}
