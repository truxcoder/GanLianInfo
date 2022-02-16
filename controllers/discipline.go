package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Disciplines struct {
	models.Discipline
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
	DictName       string `json:"dictName"`
}

func DisciplineList(c *gin.Context) {
	var mos []Disciplines
	var mo Disciplines
	selectStr := "disciplines.*,per.name as personnel_name, per.police_code as police_code, " +
		"departments.name as organ_name, departments.short_name as organ_short_name, " +
		"dis_dicts.name as dict_name "
	joinStr := "left join personnels as per on disciplines.personnel_id = per.id " +
		"left join departments on departments.id = ( select organ_id from personnels where personnels.id =  disciplines.personnel_id) " +
		"left join dis_dicts on dis_dicts.id = disciplines.dict_id"
	getList(c, "disciplines", &mo, &mos, &selectStr, &joinStr)
}

func DisciplineDetail(c *gin.Context) {
	var mos []struct {
		models.Discipline
		DictName string `json:"dictName"`
	}
	selectStr := "disciplines.*, dis_dicts.name as dict_name "
	joinStr := "left join dis_dicts on dis_dicts.id = disciplines.dict_id"
	getDetail(c, "disciplines", &mos, &selectStr, &joinStr)
}

func DisDictList(c *gin.Context) {
	var r gin.H
	var dd []models.DisDict
	result := db.Select("id", "name", "category", "term").Find(&dd)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &dd}
	}
	c.JSON(200, r)
}
