package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
	"gorm.io/gorm"
)

type PerDept struct {
	models.Personnel
	DepartmentName      string `json:"departmentName"`
	DepartmentShortName string `json:"departmentShortName"`
	OrganName           string `json:"organName"`
	OrganShortName      string `json:"organShortName"`
	Resume              string `json:"resume"`
}

type PerName struct {
	ID           int64  `json:"id,string"`
	DepartmentId string `json:"departmentId"`
	Name         string `json:"name"`
	PoliceCode   string `json:"policeCode"`
}

func PersonnelList(c *gin.Context) {
	currenPage := c.Query("currentPage")
	pageSize := c.Query("pageSize")
	var (
		pd        []PerDept
		sm        SearchMod
		r         gin.H
		err       error
		result    *gorm.DB
		count     int64         //总记录数
		paramList []interface{} //where语句参数列表
		whereStr  string        //where语句
	)

	sort := `(case when length(d.level_code)>=3 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,3)) else null end) desc,
(case when length(d.level_code)>=6 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,6)) else null end) desc,
(case when length(d.level_code)>=9 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,9)) else null end) desc,
(case when length(d.level_code)>=12 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,12)) else null end) desc,
(case when length(d.level_code)>=15 then (select ii.sort from departments ii where ii.level_code = substring(d.level_code,1,15)) else null end) desc, 
personnels.sort desc nulls first`

	selectStr := "personnels.*,d.name as department_name,d.short_name as department_short_name," +
		"o.name as organ_name,o.short_name as organ_short_name"
	joinStr := "left join departments as d on personnels.department_id = d.id " +
		"left join departments as o on personnels.organ_id = o.id"

	if err = c.BindJSON(&sm); err != nil {
		log.Error(err)
		//r = Errors.ServerError
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	whereStr, paramList = makeWhere(&sm)

	page, _ := strconv.Atoi(currenPage)
	size, _ := strconv.Atoi(pageSize)
	offset := (page - 1) * size
	//先查询数据总量并返回到前端
	err = db.Model(&models.Personnel{}).Where(whereStr, paramList...).Count(&count).Error
	if err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	if count == 0 {
		r = GetResponse(ResNoData)
		c.JSON(200, r)
		return
	}
	result = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where(whereStr, paramList...).Order(sort).Limit(size).Offset(offset).Find(&pd)

	if result.Error != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &pd, "count": count}
	c.JSON(200, r)
}

func PersonnelDetail(c *gin.Context) {
	var r gin.H
	var p models.Personnel
	var pd PerDept
	//var posts []models.Post
	var trains []models.Training
	selectStr := "personnels.*,d1.name as department_name,d1.short_name as department_short_name," +
		"d2.name as organ_name,d2.short_name as organ_short_name," +
		"resumes.content as resume"
	joinStr := "left join departments as d1 on personnels.department_id = d1.id " +
		"left join departments as d2 on personnels.organ_id = d2.id " +
		"left join resumes on resumes.personnel_id = personnels.id"
	var id ID
	if err := c.ShouldBindJSON(&id); err != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	//db.Debug().Preload("Department").Preload("Posts").First(&p, "id = ?", personnelId)
	db.Model(&p).Select(selectStr).Joins(joinStr).Where("personnels.id = ?", id.ID).First(&pd)
	//db.Where("personnel_id = ?", id.ID).Find(&posts)
	db.Table("trainings").Where("id in (?)", db.Table("person_trains").Select("train_id").Where("personnel_id = ?", id.ID)).Find(&trains)
	//r = gin.H{"code": 20000, "data": &pd, "posts": &posts, "trains": &trains}
	r = gin.H{"code": 20000, "data": &pd, "trains": &trains}
	c.JSON(200, r)
}

// PersonnelBaseList 基本信息列表，用于人员选择等
func PersonnelBaseList(c *gin.Context) {
	var p []PerDept
	var err error
	var r gin.H
	var result *gorm.DB
	var mo struct {
		AccountId string `json:"accountId"`
		OrganId   string `json:"organId"`
	}
	if c.ShouldBindJSON(&mo) != nil {
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	canGlobal, _ := enforcer.Enforce(mo.AccountId, "Personnel", "GLOBAL")
	selectStr := "personnels.*, departments.name as organ_name,departments.short_name as organ_short_name"
	joinStr := "left join departments on personnels.organ_id = departments.id"
	if canGlobal {
		result = db.Table("personnels").Select(selectStr).Joins(joinStr).Where("status = 1").Find(&p)
	} else {
		result = db.Table("personnels").Select(selectStr).Joins(joinStr).Where("status = 1 AND organ_id = ?", mo.OrganId).Find(&p)
	}
	err = result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
	} else if result.RowsAffected == 0 {
		//r = Errors.NoData
		r = GetError(CodeNoData)
	} else {
		r = gin.H{"code": 20000, "data": &p}
	}
	c.JSON(200, r)
}

// PersonnelNameList 人员姓名列表，用于前端穿梭框选择
func PersonnelNameList(c *gin.Context) {
	var p []PerName
	var err error
	var r gin.H
	var result *gorm.DB
	var mo struct {
		AccountId string `json:"accountId"`
		OrganId   string `json:"organId"`
	}
	if c.ShouldBindJSON(&mo) != nil {
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	canGlobal, _ := enforcer.Enforce(mo.AccountId, "Personnel", "GLOBAL")
	selectStr := "id, name, police_code, department_id"
	if canGlobal {
		result = db.Table("personnels").Select(selectStr).Where("status = 1").Find(&p)
	} else {
		result = db.Table("personnels").Select(selectStr).Where("status = 1 AND organ_id = ?", &mo.OrganId).Find(&p)
	}
	err = result.Error
	if err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	if result.RowsAffected == 0 {
		r = GetError(CodeNoData)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &p}
	c.JSON(200, r)
}

func PersonnelUpdate(c *gin.Context) {
	var p models.Personnel
	var r gin.H
	if c.ShouldBindJSON(&p) != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
	} else {
		p.Birthday = utils.GetBirthdayFromIdCode(p.IdCode)
		db.Model(&p).Updates(&p)
		r = gin.H{"message": "更新成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func PersonnelDelete(c *gin.Context) {
	var (
		err error
		r   gin.H
		id  int64
	)

	if id, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	result := db.Delete(models.Personnel{}, id)
	err = result.Error
	if err != nil {
		log.Error(err)
		r = GetError(CodeServer)
	} else {
		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
		r = gin.H{"message": message, "code": 20000}
	}
	c.JSON(200, r)
	return
}

//func PersonnelDelete(c *gin.Context) {
//	var err error
//	var id IdStruct
//	var r gin.H
//	if err = c.ShouldBindJSON(&id); err != nil {
//		r = GetError(CodeBind)
//		log.Error(err)
//		c.JSON(200, r)
//		return
//	}
//	result := db.Delete(models.Personnel{}, &id.Id)
//	err = result.Error
//	if err != nil {
//		log.Error(err)
//		r = GetError(CodeServer)
//	} else {
//		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
//		r = gin.H{"message": message, "code": 20000}
//	}
//	c.JSON(200, r)
//	return
//}

func PersonnelResume(c *gin.Context) {
	var resume, res models.Resume
	var r gin.H
	if c.ShouldBindJSON(&resume) != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	result := db.Table("resumes").Where("personnel_id = ?", resume.PersonnelId).Limit(1).Find(&res)
	if result.RowsAffected > 0 {
		db.Table("resumes").Where("id = ?", res.ID).Updates(&resume)
	} else {
		db.Table("resumes").Create(&resume)
	}
	r = gin.H{"message": "更新成功！", "code": 20000}
	c.JSON(200, r)
	return
}

func UpdateIdCode(c *gin.Context) {
	var r gin.H
	var p struct {
		ID     int64  `json:"id,string"`
		IdCode string `json:"idCode"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	result := db.Table("personnels").Where("id = ?", p.ID).Update("id_code", p.IdCode)
	if result.Error != nil || result.RowsAffected == 0 {
		//r = Errors.Update
		r = GetError(CodeUpdate)
		log.Error(result.Error)
		c.JSON(200, r)
		return
	}
	setIdCodeMap()
	r = gin.H{"code": 20000, "message": "身份证号码更新成功!"}
	c.JSON(200, r)
}

func EduDictList(c *gin.Context) {
	var dicts []models.EduDict
	var r gin.H
	result := db.Order("sort asc").Find(&dicts)
	err := result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		log.Error(err)
	} else {
		r = gin.H{"code": 20000, "data": &dicts}
	}
	//time.Sleep(4 * time.Second)
	c.JSON(200, r)
}

// GetPersonOrganId 获取人员单位ID，用在人员detail页面权限认证
func GetPersonOrganId(c *gin.Context) {
	var id struct {
		ID int64 `json:"id,string"`
	}
	var r gin.H
	if c.ShouldBindJSON(&id) != nil {
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	var u struct {
		OrganId string
	}

	if id.ID == 0 {
		r = gin.H{"code": 20000, "data": ""}
		c.JSON(200, r)
		return
	}
	db.Table("personnels").Select("organ_id").Where("id = ?", id.ID).Limit(1).Find(&u)
	r = gin.H{"code": 20000, "data": u.OrganId}
	c.JSON(200, r)
	return
}

// GetPersonOrgans 获取所有人员的id
func GetPersonOrgans(c *gin.Context) {
	var r gin.H
	var u []struct {
		ID      int64 `json:"id,string"`
		OrganId string
	}
	_map := make(map[string]string)
	db.Table("personnels").Select("id, organ_id").Find(&u)
	for _, v := range u {
		_map[strconv.FormatInt(v.ID, 10)] = v.OrganId
	}
	r = gin.H{"code": 20000, "data": _map}
	c.JSON(200, r)
	return
}
