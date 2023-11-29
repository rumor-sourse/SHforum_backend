package logic

import (
	"SHforum_backend/dao/mysql"
	"SHforum_backend/models/response"
)

func GetCommunityList() (data []*response.CommunityResponse, err error) {
	communities, err := mysql.GetCommunityList()
	if err != nil {
		return nil, err
	}
	for _, community := range communities {
		data = append(data, &response.CommunityResponse{
			ID:   community.ID,
			Name: community.Name,
		})
	}
	return
}

func GetCommunityDetail(id int64) (data *response.CommunityDetailResponse, err error) {
	community, err := mysql.GetCommunityDetailByID(id)
	if err != nil {
		return nil, err
	}
	data = &response.CommunityDetailResponse{
		ID:           community.ID,
		Name:         community.Name,
		Introduction: community.Introduction,
	}
	return
}
