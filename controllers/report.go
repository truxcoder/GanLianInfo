package controllers

import (
	"GanLianInfo/models"

	"github.com/Insua/gorm-dm8/datatype"

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
	result := db.Table("reports").Omit("reports.intro, reports.steps").Find(&mos)
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
		r = gin.H{"code": 20000, "data": &mos}
	}
	c.JSON(200, r)
}

// ReportOne 一条记录的额外信息
func ReportOne(c *gin.Context) {
	var err error
	var r gin.H
	var id struct {
		ID int64 `json:"id"`
	}
	var mo struct {
		Intro datatype.Clob `json:"intro"`
		Steps datatype.Clob `json:"steps"`
	}
	if err = c.ShouldBindJSON(&id); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Table("reports").Select("intro, steps").Where("id = ?", id.ID).Limit(1).Find(&mo)
	r = gin.H{"code": 20000, "data": &mo}
	c.JSON(200, r)
}

func ReportDetail(c *gin.Context) {
	var err error
	var mos []Report
	var r gin.H
	var id struct {
		ID string `json:"id"`
	}
	if err = c.ShouldBindJSON(&id); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Table("reports").Omit("reports.intro, reports.steps").Where("id in (?)", db.Table("person_reports").Select("report_id").Where("personnel_id = ?", id.ID)).Find(&mos)
	selectStr := "personnels.id, personnels.name, personnels.police_code, personnels.organ_id, d.name as organ_name, d.short_name as organ_short_name"
	joinStr := "left join departments as d on personnels.organ_id = d.id"
	for i := 0; i < len(mos); i++ {
		var personnels []*PersonSimple
		db.Table("personnels").Select(selectStr).Joins(joinStr).Where("personnels.id in (?)", db.Table("person_reports").Select("personnel_id").Where("report_id = ?", mos[i].ID)).Find(&personnels)
		mos[i].Personnels = append(mos[i].Personnels, personnels...)
	}
	r = gin.H{"code": 20000, "data": &mos}
	c.JSON(200, r)
}

func ReportSteps(c *gin.Context) {
	var r gin.H
	var err error
	var mo struct {
		ID    int64         `json:"id"`
		Steps datatype.Clob `json:"steps"`
	}
	if err = c.ShouldBindJSON(&mo); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Table("reports").Select("steps").Where("id = ?", &mo.ID).Limit(1).Find(&mo)
	r = gin.H{"code": 20000, "data": &mo}
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
	if len(model.Person) == 0 {
		r = gin.H{"message": "添加成功！", "code": 20000}
		c.JSON(200, r)
		return
	}
	id := model.Report.ID
	for _, v := range model.Person {
		personReports = append(personReports, models.PersonReport{ReportId: id, PersonnelId: v})
	}
	db.Table("person_reports").Create(personReports)
	r = gin.H{"message": "添加成功！", "code": 20000}
	c.JSON(200, r)
}

func ReportUpdate(c *gin.Context) {
	var r gin.H
	var err error
	var personReports []models.PersonReport
	var model struct {
		Report models.Report `json:"report"`
		Add    []string      `json:"add"`
		Del    []string      `json:"del"`
	}
	if err = c.ShouldBindJSON(&model); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	db.Model(&model.Report).Updates(&model.Report)
	id := model.Report.ID
	if len(model.Add) > 0 {
		for _, v := range model.Add {
			personReports = append(personReports, models.PersonReport{ReportId: id, PersonnelId: v})
		}
		db.Table("person_reports").Create(personReports)
	}
	if len(model.Del) > 0 {
		result := db.Where("report_id = ? and personnel_id in ?", id, model.Del).Delete(models.PersonReport{})
		err = result.Error
		if err != nil {
			log.Error(err)
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
	}
	r = gin.H{"message": "修改成功！", "code": 20000}
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
