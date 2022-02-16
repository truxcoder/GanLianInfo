package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Punishes struct {
	models.Punish
	PersonnelName string `json:"personnelName"`
	PoliceCode    string `json:"policeCode"`
}

func PunishList(c *gin.Context) {
	var mos []Punishes
	var mo Punishes
	selectStr := "punishes.*,per.name as personnel_name, per.police_code as police_code "
	joinStr := "left join personnels as per on punishes.personnel_id = per.id "
	getList(c, "punishes", &mo, &mos, &selectStr, &joinStr)
}

func PunishDetail(c *gin.Context) {
	var mos []models.Punish
	var selectStr string
	var joinStr string
	getDetail(c, "punishes", &mos, &selectStr, &joinStr)
}
