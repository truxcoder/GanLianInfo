package controllers

import (
	"GanLianInfo/models"

	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

type Talent struct {
	models.Talent
	PersonnelName  string `json:"personnelName"`
	Birthday       string `json:"birthday"`
	PoliceCode     string `json:"policeCode"`
	Gender         string `json:"gender"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func TalentList(c *gin.Context) {
	var mos []Talent
	var mo Talent
	selectStr := "talents.*,per.name as personnel_name, per.police_code as police_code, per.birthday as birthday, per.gender as gender, " +
		"departments.name as organ_name, departments.short_name as organ_short_name"
	joinStr := "left join personnels as per on per.id = talents.personnel_id " +
		"left join departments on departments.id = ( select organ_id from personnels where personnels.id = talents.personnel_id)"
	getList(c, "talents", &mo, &mos, &selectStr, &joinStr)
}

func TalentAdd(c *gin.Context) {
	var (
		err   error
		r     gin.H
		mo    models.Talent
		count int64
	)

	if err = c.ShouldBindJSON(&mo); err != nil {
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}
	if db.Table("talents").Where("personnel_id = ? AND category = ?", mo.PersonnelId, mo.Category).Count(&count); count > 0 {
		log.Successf("count:%d\n", count)
		r = gin.H{"code": 503, "message": "系统内已存在该人员信息，请勿重复添加！"}
		c.JSON(200, r)
		return
	}
	db.Create(&mo)
	r = gin.H{"code": 20000, "message": "添加成功！"}
	c.JSON(200, r)
}
