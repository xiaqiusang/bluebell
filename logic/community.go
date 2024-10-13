package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunitylist() ([]*models.Community, error) {
	//查询所有的community并返回
	return mysql.GetCommunitylist()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailbyID(id)
}
