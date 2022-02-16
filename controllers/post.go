package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Posts struct {
	models.Post
	PersonnelName string `json:"personnelName"`
	PoliceCode    string `json:"policeCode"`
	PositionName  string `json:"positionName"`
	LevelName     string `json:"levelName"`
}

type PostDetailStruct struct {
	models.Post
	PositionName string `json:"positionName"`
	LevelName    string `json:"levelName"`
}

func PostList(c *gin.Context) {
	var mos []Posts
	var mo Posts
	selectStr := "posts.*,per.name as personnel_name, per.police_code as police_code," +
		"p.name as position_name, l.name as level_name"
	joinStr := "left join personnels as per on posts.personnel_id = per.id " +
		"left join positions as p on posts.position_id = p.id " +
		"left join levels as l on posts.level_id = l.id "
	getList(c, "posts", &mo, &mos, &selectStr, &joinStr)
}

func PostDetail(c *gin.Context) {
	var mos []PostDetailStruct
	selectStr := "posts.*,p.name as position_name, l.name as level_name"
	joinStr := "left join positions as p on posts.position_id = p.id " +
		"left join levels as l on posts.level_id = l.id"
	getDetail(c, "posts", &mos, &selectStr, &joinStr)
}
