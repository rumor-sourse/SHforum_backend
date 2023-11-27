package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/models"
)

func GetCommunityList() (data []models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailByID(id)
}
