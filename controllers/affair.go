package controllers

import (
	"GanLianInfo/models"

	"github.com/Insua/gorm-dm8/datatype"
	"github.com/gin-gonic/gin"
)

type Affair struct {
	models.Affair
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func AffairList(c *gin.Context) {
	var mos []Affair
	var mo Affair
	selectStr := "affairs.id,affairs.personnel_id,affairs.title,affairs.category,personnels.name as personnel_name, personnels.police_code as police_code" +
		", departments.name as organ_name, departments.short_name as organ_short_name"
	joinStr := "left join personnels on affairs.personnel_id = personnels.id " +
		"left join departments on departments.id = (select organ_id from personnels where personnels.id = affairs.personnel_id)"
	getList(c, "affairs", &mo, &mos, &selectStr, &joinStr)
}

func AffairDetail(c *gin.Context) {
	var mos []models.Affair
	var selectStr string
	var joinStr string
	getDetail(c, "affairs", &mos, &selectStr, &joinStr)
}

func AffairOne(c *gin.Context) {
	var err error
	var r gin.H
	var id ID
	var mo struct {
		Intro datatype.Clob `json:"intro"`
	}
	if err = c.ShouldBindJSON(&id); err != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	db.Table("affairs").Select("intro").Where("id = ?", id.ID).Limit(1).Find(&mo)
	r = gin.H{"code": 20000, "data": &mo}
	c.JSON(200, r)
}
