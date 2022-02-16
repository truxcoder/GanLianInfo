package controllers

import (
	"GanLianInfo/models"

	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

type PersonSimple struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PoliceCode     string `json:"policeCode"`
	OrganId        string `json:"organId"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

type Report struct {
	models.Report
	Personnels []*PersonSimple `json:"personnels" gorm:"-"`
}

func ReportList(c *gin.Context) {
	var mos []Report
	var pr []models.PersonReport
	var r gin.H
	var ids []int64
	var personnels []PersonSimple
	var personMap = make(map[string]*PersonSimple)
	var personReportMap = make(map[int64][]string)
	result := db.Table("reports").Omit("reports.personnels").Find(&mos)
	for _, v := range mos {
		ids = append(ids, v.ID)
	}
	db.Table("person_reports").Find(&pr)
	selectStr := "personnels.id, personnels.name, personnels.police_code, personnels.organ_id, d.name as organ_name, d.short_name as organ_short_name"
	joinStr := "left join departments as d on personnels.organ_id = d.id"
	db.Table("personnels").Select(selectStr).Joins(joinStr).Where("personnels.id in (?)", db.Table("person_reports").Select("personnel_id").Where("report_id in ?", ids)).Find(&personnels)
	for i := 0; i < len(personnels); i++ {
		personMap[personnels[i].ID] = &personnels[i]
	}
	for _, v := range pr {
		personReportMap[v.ReportId] = append(personReportMap[v.ReportId], v.PersonnelId)
	}
	for i := 0; i < len(mos); i++ {
		personList := personReportMap[mos[i].ID]
		for _, v := range personList {
			mos[i].Personnels = append(mos[i].Personnels, personMap[v])
		}
	}
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &mos, "personnels": &personnels, "map": personReportMap, "map2": personMap}
	}
	c.JSON(200, r)
}

func ReportAdd(c *gin.Context) {
	var r gin.H
	var err error
	var personReports []models.PersonReport
	var model struct {
		Report models.Report `json:"report"`
		Person []string      `json:"person"`
	}
	if err = c.ShouldBindJSON(&model); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Table("reports").Create(&model.Report)
	id := model.Report.ID
	log.Successf("ID:%s\n", id)
	for _, v := range model.Person {
		personReports = append(personReports, models.PersonReport{ReportId: id, PersonnelId: v})
	}
	db.Table("person_reports").Create(personReports)
	r = gin.H{"message": "添加成功！", "code": 20000}
	c.JSON(200, r)
}

func PersonReportAdd(c *gin.Context) {
	var r gin.H
	var mos []models.PersonReport
	if err := c.ShouldBindJSON(&mos); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Create(&mos)
	r = gin.H{"message": "添加成功！", "code": 20000}
	c.JSON(200, r)
}
