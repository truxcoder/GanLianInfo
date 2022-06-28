package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"time"

	"github.com/Insua/gorm-dm8/datatype"
	jsoniter "github.com/json-iterator/go"

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
		r = GetError(CodeBind)
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

func TalentDetailList(c *gin.Context) {
	var (
		ids []int64
		err error
		r   gin.H
		mos []Talent
		mo  struct {
			Res string `json:"res"`
		}
		selectStr = "talents.*,per.name as personnel_name, per.police_code as police_code, per.birthday as birthday, per.gender as gender, " +
			"departments.name as organ_name, departments.short_name as organ_short_name"
		joinStr = "left join personnels as per on per.id = talents.personnel_id " +
			"left join departments on departments.id = ( select organ_id from personnels where personnels.id = talents.personnel_id)"
	)
	if err = c.BindJSON(&mo); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	if err = jsoniter.UnmarshalFromString(mo.Res, &ids); err != nil {
		log.Error(err)
		r = GetError(CodeParse)
		c.JSON(200, r)
		return
	}
	db.Table("talents").Select(selectStr).Joins(joinStr).Where("talents.id in ?", ids).Find(&mos)
	r = gin.H{"code": 20000, "data": &mos}
	c.JSON(200, r)
}

func TalentPickList(c *gin.Context) {
	var mos []struct {
		models.TalentPick
		PickerName string `json:"pickerName"`
	}
	var mo models.TalentPick
	selectStr := "talent_picks.*,per.name as picker_name"
	joinStr := "left join personnels as per on per.id = talent_picks.picker_id "
	getList(c, "talent_picks", &mo, &mos, &selectStr, &joinStr)
}

func TalentPickAdd(c *gin.Context) {
	var (
		ids []int64
		err error
		r   gin.H
		res string
		mo  struct {
			Category int8   `json:"category"`
			Title    string `json:"title"`
			PickerId int64  `json:"pickerId,string"`
			Total    int    `json:"total"`
		}
	)
	if err = c.BindJSON(&mo); err != nil {
		log.Error(err)
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	db.Model(&models.Talent{}).Where("category = ?", &mo.Category).Pluck("id", &ids)
	ids = utils.GetRandIdList(ids, mo.Total)
	if res, err = jsoniter.MarshalToString(ids); err != nil {
		log.Error(err)
		r = GetError(CodeParse)
		c.JSON(200, r)
		return
	}
	tp := models.TalentPick{PickerId: mo.PickerId, Category: mo.Category, Title: mo.Title, PickDate: time.Now(), Res: datatype.Clob(res)}
	db.Create(&tp)
	r = gin.H{"message": "抽取成功!", "code": 20000}
	c.JSON(200, r)
}
