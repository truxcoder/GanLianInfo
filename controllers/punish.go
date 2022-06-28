package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Punishes struct {
	models.Punish
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func PunishList(c *gin.Context) {
	var mos []Punishes
	var mo struct {
		models.Punish
		Intercept string `json:"intercept" gorm:"-" sql:""`
	}
	selectStr := "punishes.*,per.name as personnel_name, per.police_code as police_code," +
		"departments.name as organ_name, departments.short_name as organ_short_name "
	joinStr := "left join personnels as per on punishes.personnel_id = per.id " +
		"left join departments on departments.id = per.organ_id "
	getList(c, "punishes", &mo, &mos, &selectStr, &joinStr)
}

func PunishDetail(c *gin.Context) {
	var mos []models.Punish
	var selectStr string
	var joinStr string
	getDetail(c, "punishes", &mos, &selectStr, &joinStr)
}
