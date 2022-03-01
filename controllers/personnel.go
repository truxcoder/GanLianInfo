package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"strconv"
	"time"

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
	ID           string `json:"id"`
	DepartmentId string `json:"departmentId"`
	Name         string `json:"name"`
	PoliceCode   string `json:"policeCode"`
}

type SearchMod struct {
	AgeStart                int8     `json:"ageStart"`
	AgeEnd                  int8     `json:"ageEnd"`
	Gender                  string   `json:"gender"`
	Nation                  []string `json:"nation"`
	Political               []string `json:"political"`
	FullTimeEdu             []string `json:"fullTimeEdu"`
	FullTimeMajor           []string `json:"fullTimeMajor"`
	PartTimeEdu             []string `json:"partTimeEdu"`
	OrganID                 string   `json:"organId"`
	ProCert                 []string `json:"proCert"`
	IsSecret                string   `json:"isSecret"`
	PassExamDay             string   `json:"passExamDay"`
	HasPassport             string   `json:"hasPassport"`
	HasAppraisalIncompetent string   `json:"hasAppraisalIncompetent"`
	HasReport               string   `json:"hasReport"`
	Award                   []int8   `json:"award"`
	Punish                  []int8   `json:"punish"`
}

func PersonnelList(c *gin.Context) {
	queryMeans := c.Query("queryMeans") //请求方式，是前端分页还是后端分页
	currenPage := c.Query("currentPage")
	pageSize := c.Query("pageSize")
	searchMeans := c.Query("searchMeans") //查询类型，是高级查询还是普通查询
	var pd []PerDept
	var sm SearchMod
	var personnel models.Personnel
	var r gin.H
	var err error
	var result *gorm.DB
	var count int64             //总记录数
	var paramList []interface{} //where语句参数列表
	var whereStr string         //where语句
	whereName := "1 = 1"        //简易搜索时名称用模糊查询
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

	//先判断搜索方式是高级搜索还是普通搜索，分别绑定不同的实例，后面所有都要根据搜索方式的不同利用不同的where语句
	if searchMeans == "advance" {
		if err = c.BindJSON(&sm); err != nil {
			log.Error(err)
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
		whereStr, paramList = makeWhere(&sm)
	} else {
		if err = c.BindJSON(&personnel); err != nil {
			log.Error(err)
			r = Errors.ServerError
			c.JSON(200, r)
			return
		}
		if personnel.Name != "" {
			whereName = "personnels.name like '%" + personnel.Name + "%'"
			personnel.Name = ""
		}
	}
	//后端分页情况
	if queryMeans == "backend" {
		page, _ := strconv.Atoi(currenPage)
		size, _ := strconv.Atoi(pageSize)
		offset := (page - 1) * size
		//先查询数据总量并返回到前端
		if searchMeans == "advance" {
			err = db.Model(&models.Personnel{}).Where(whereStr, paramList...).Count(&count).Error
		} else {
			err = db.Model(&models.Personnel{}).Where(&personnel).Where(whereName).Count(&count).Error
		}
		if err != nil {
			r = Errors.ServerError
		} else if count != 0 {
			//result = db.Preload("Organ").Where(&personnel).Limit(size).Offset(offset).Find(&p)
			if searchMeans == "advance" {
				result = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where(whereStr, paramList...).Order(sort).Limit(size).Offset(offset).Find(&pd)
			} else {
				result = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where(&personnel).Where(whereName).Order(sort).Limit(size).Offset(offset).Find(&pd)
			}
			if result.Error != nil {
				r = Errors.ServerError
			} else {
				r = gin.H{"code": 20000, "data": &pd, "count": count}
			}
		} else {
			r = Errors.NoData
		}
	} else { //前端分页情况
		//result = db.Preload("Organ").Where(&personnel).Find(&p)
		if searchMeans == "advance" {
			result = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where(whereStr, paramList...).Order(sort).Find(&pd)
		} else {
			result = db.Model(&models.Personnel{}).Select(selectStr).Joins(joinStr).Where(&personnel).Order(sort).Find(&pd)
		}
		err = result.Error
		if err != nil {
			r = Errors.ServerError
		} else if result.RowsAffected == 0 {
			r = Errors.NoData
		} else {
			r = gin.H{"code": 20000, "data": &pd, "count": result.RowsAffected}
		}
	}
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
	var id struct {
		ID string `json:"id"`
	}
	if c.ShouldBindJSON(&id) != nil {
		r = Errors.ServerError
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

func GetPersonnelName(c *gin.Context) {
	var p []PerDept
	var err error
	var r gin.H
	var result *gorm.DB
	var organ struct {
		PersonnelId string `json:"personnelId"`
		OrganId     string `json:"organId"`
	}
	if c.ShouldBindJSON(&organ) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	canGlobal, _ := enforcer.Enforce(organ.PersonnelId, "Personnel", "GLOBAL")
	selectStr := "personnels.*, departments.name as organ_name,departments.short_name as organ_short_name"
	joinStr := "left join departments on personnels.organ_id = departments.id"
	if canGlobal {
		result = db.Table("personnels").Select(selectStr).Joins(joinStr).Find(&p)
	} else {
		result = db.Table("personnels").Select(selectStr).Joins(joinStr).Where("organ_id = ?", organ.OrganId).Find(&p)
	}
	err = result.Error
	if err != nil {
		r = Errors.ServerError
	} else if result.RowsAffected == 0 {
		r = Errors.NoData
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
	var organ struct {
		PersonnelId string `json:"personnelId"`
		OrganId     string `json:"organId"`
	}
	if c.ShouldBindJSON(&organ) != nil {
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	canGlobal, _ := enforcer.Enforce(organ.PersonnelId, "Personnel", "GLOBAL")
	selectStr := "id, name, police_code, department_id"
	if canGlobal {
		result = db.Table("personnels").Select(selectStr).Find(&p)
	} else {
		result = db.Table("personnels").Select(selectStr).Where("organ_id = ?", organ.OrganId).Find(&p)
	}
	err = result.Error
	if err != nil {
		r = Errors.ServerError
	} else if result.RowsAffected == 0 {
		r = Errors.NoData
	} else {
		r = gin.H{"code": 20000, "data": &p}
	}
	c.JSON(200, r)
}

func SearchPersonnelName(c *gin.Context) {
	var p []models.Personnel
	var name struct {
		Name string `json:"name"`
	}
	var err error
	var r gin.H
	var result *gorm.DB
	if err = c.BindJSON(&name); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	result = db.Select("id", "name", "police_code").Where("name LIKE ?", name.Name+"%").Find(&p)
	err = result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &p}
	}
	c.JSON(200, r)
}

func PersonnelUpdate(c *gin.Context) {
	var p models.Personnel
	var r gin.H
	if c.ShouldBindJSON(&p) != nil {
		r = Errors.ServerError
	} else {
		p.Birthday = utils.GetBirthdayFromIdCode(p.IdCode)
		db.Model(&p).Updates(&p)
		r = gin.H{"message": "更新成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func PersonnelResume(c *gin.Context) {
	var resume, res models.Resume
	var r gin.H
	if c.ShouldBindJSON(&resume) != nil {
		r = Errors.ServerError
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

func EduDictList(c *gin.Context) {
	var dicts []models.EduDict
	var r gin.H
	result := db.Order("sort asc").Find(&dicts)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &dicts}
	}
	//time.Sleep(4 * time.Second)
	c.JSON(200, r)
}

func yearsAgo(year int8) time.Time {
	now := time.Now()
	return now.AddDate(int(0-year), 0, 0)
}

func makeWhere(sm *SearchMod) (string, []interface{}) {
	var paramList []interface{}
	whereStr := " 1 = 1 "
	if sm.AgeStart != 0 && sm.AgeEnd != 0 {
		whereStr += " AND birthday >= ? AND birthday <= ?"
		paramList = append(paramList, yearsAgo(sm.AgeEnd), yearsAgo(sm.AgeStart))
	} else if sm.AgeStart != 0 {
		whereStr += " AND birthday <= ?"
		paramList = append(paramList, yearsAgo(sm.AgeStart))
	} else if sm.AgeEnd != 0 {
		whereStr += " AND birthday >= ?"
		paramList = append(paramList, yearsAgo(sm.AgeEnd))
	}
	if sm.Gender != "" {
		whereStr += " AND gender = ?"
		paramList = append(paramList, sm.Gender)
	}
	if len(sm.Nation) > 0 {
		hanFound := false
		otherFound := false
		for _, v := range sm.Nation {
			if v == "少数民族" {
				otherFound = true
			}
			if v == "汉族" {
				hanFound = true
			}
			if hanFound && otherFound {
				break
			}
		}
		if otherFound && !hanFound {
			whereStr += " AND nation <> ?"
			paramList = append(paramList, "汉族")
		} else if !otherFound {
			whereStr += " AND nation in ?"
			paramList = append(paramList, sm.Nation)
		}
	}
	if len(sm.Political) > 0 {
		whereStr += " AND political in ?"
		paramList = append(paramList, sm.Political)
	}
	if len(sm.FullTimeEdu) > 0 {
		whereStr += " AND full_time_edu in ?"
		paramList = append(paramList, sm.FullTimeEdu)
	}
	if len(sm.FullTimeMajor) > 0 {
		//whereStr += " AND full_time_major in ?"
		//paramList = append(paramList, sm.FullTimeMajor)
		for k, v := range sm.FullTimeMajor {
			if k == 0 {
				whereStr += " AND (full_time_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			} else {
				whereStr += " OR full_time_major LIKE ?"
				paramList = append(paramList, "%"+v+"%")
			}
		}
		whereStr += ")"
	}
	if len(sm.PartTimeEdu) > 0 {
		whereStr += " AND part_time_edu in ?"
		paramList = append(paramList, sm.PartTimeEdu)
	}
	if sm.OrganID != "" {
		whereStr += " AND organ_id = ?"
		paramList = append(paramList, sm.OrganID)
	}
	if len(sm.ProCert) > 0 {
		//whereStr += " AND pro_cert in ?"
		//paramList = append(paramList, sm.ProCert)
		for _, v := range sm.ProCert {
			whereStr += " AND pro_cert like ?"
			paramList = append(paramList, "%"+v+"%")
		}
	}
	if sm.IsSecret != "" {
		whereStr += " AND is_secret = ?"
		paramList = append(paramList, getBool(sm.IsSecret))
	}
	if sm.PassExamDay != "" {
		if sm.PassExamDay == "是" {
			whereStr += " AND pass_exam_day > ?"
		} else {
			whereStr += " AND (pass_exam_day < ? or pass_exam_day is null)"
		}
		paramList = append(paramList, yearsAgo(3))
	}
	if sm.HasPassport != "" {
		//whereStr += " AND has_passport = ?"
		//paramList = append(paramList, getBool(sm.HasPassport))
		if sm.HasPassport == "是" {
			whereStr += " AND (passport is not null and json_value(personnels.passport, '$[0]' RETURNING number) <> 0)"
		} else {
			whereStr += " AND (passport is null or json_value(personnels.passport, '$[0]' RETURNING number) = 0)"
		}
	}
	if sm.HasAppraisalIncompetent != "" {
		if sm.HasAppraisalIncompetent == "是" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("appraisals").Select("personnel_id").Where("conclusion in ('不称职','不确定等次')"))
	}
	if sm.HasReport != "" {
		if sm.HasReport == "是" {
			whereStr += " AND personnels.id in (?)"
		} else {
			whereStr += " AND personnels.id not in (?)"
		}
		paramList = append(paramList, db.Table("person_reports").Select("personnel_id").Where("report_id in (?)",
			db.Table("reports").Select("id").Where("step <> ?", 99)))
	}
	if len(sm.Award) > 0 {
		whereStr += " AND personnels.id in (?)"
		paramList = append(paramList, db.Table("awards").Select("personnel_id").Where("grade in ?", sm.Award))
	}
	if len(sm.Punish) > 0 {
		whereStr += " AND personnels.id in (?)"
		paramList = append(paramList, db.Table("punishes").Select("personnel_id").Where("grade in ?", sm.Punish))
	}
	return whereStr, paramList
}

func getBool(str string) int {
	if str == "是" {
		return 1
	}
	return 0
}
