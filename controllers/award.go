package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Awards struct {
	models.Award
	PersonnelName string `json:"personnelName"`
	PoliceCode    string `json:"policeCode"`
}

func AwardList(c *gin.Context) {
	var mos []Awards
	var mo Awards
	selectStr := "awards.*,per.name as personnel_name, per.police_code as police_code "
	joinStr := "left join personnels as per on awards.personnel_id = per.id "
	getList(c, "awards", &mo, &mos, &selectStr, &joinStr)
}

func AwardDetail(c *gin.Context) {
	var mos []models.Award
	var selectStr string
	var joinStr string
	getDetail(c, "awards", &mos, &selectStr, &joinStr)
}
