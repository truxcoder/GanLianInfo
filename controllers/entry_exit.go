package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type EntryExit struct {
	models.EntryExit
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func EntryExitList(c *gin.Context) {
	var mos []EntryExit
	var mo EntryExit
	selectStr := "entry_exits.*,per.name as personnel_name, per.police_code as police_code," +
		"departments.name as organ_name, departments.short_name as organ_short_name"
	joinStr := "left join personnels as per on entry_exits.personnel_id = per.id " +
		"left join departments on departments.id = ( select organ_id from personnels where personnels.id = entry_exits.personnel_id) "
	getList(c, "entry_exits", &mo, &mos, &selectStr, &joinStr)
}

//func EntryExitDetail(c *gin.Context) {
//	var mos []PostDetailStruct
//	selectStr := "posts.*,p.name as position_name, l.name as level_name"
//	joinStr := "left join positions as p on posts.position_id = p.id " +
//		"left join levels as l on posts.level_id = l.id"
//	getDetail(c, "posts", &mos, &selectStr, &joinStr)
//}
