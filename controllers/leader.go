package controllers

import (
	"GanLianInfo/models"
	"github.com/gin-gonic/gin"
)

// LeaderList 获取班子成员列表数据
func LeaderList(c *gin.Context) {
	var mos []LeaderTeam
	var mo models.Leader
	selectStr := "leaders.*,per.name as personnel_name, per.police_code as police_code," +
		"d.name as organ_name, d.short_name as organ_short_name"
	joinStr := "left join personnels as per on leaders.personnel_id = per.id " +
		"left join departments as d on leaders.organ_id = d.id "
	getList(c, "leaders", &mo, &mos, &selectStr, &joinStr)
}
