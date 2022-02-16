package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type AppPerson struct {
	models.Appraisal
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
}

type AppOrgan struct {
	models.Appraisal
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func AppraisalList(c *gin.Context) {
	var mos []AppPerson
	var mo AppPerson
	selectStr := "appraisals.*,d.name as organ_name,d.short_name as organ_short_name," +
		"p.name as personnel_name,p.police_code as police_code"
	joinStr := "left join departments as d on appraisals.organ_id = d.id " +
		"left join personnels as p on appraisals.personnel_id = p.id"
	getList(c, "appraisals", &mo, &mos, &selectStr, &joinStr)

}

func AppraisalDetail(c *gin.Context) {
	var mos []AppOrgan
	selectStr := "appraisals.*,d.name as organ_name,d.short_name as organ_short_name"
	joinStr := "left join departments as d on appraisals.organ_id = d.id "
	getDetail(c, "appraisals", &mos, &selectStr, &joinStr)
}
