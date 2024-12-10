package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

type Organ struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Sort      int    `json:"sort"`
}

type Headcount struct {
	models.Department
	Use int `json:"use"`
}

func DepartmentList(c *gin.Context) {
	var r gin.H
	d, err := getDepartmentSlice()
	if err != nil {
		//r = Errors.DatabaseError
		r = GetError(CodeDatabase)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &d}
	c.JSON(200, r)
}

func OrganList(c *gin.Context) {
	var o []Organ
	var r gin.H
	result := db.Model(&models.Department{}).Where("dept_type = ?", 1).Order("sort asc").Find(&o)
	err := result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
	} else {
		r = gin.H{"code": 20000, "data": &o}
	}
	c.JSON(200, r)
}

// HeadcountList 编制数列表
func HeadcountList(c *gin.Context) {
	var (
		h     []Headcount
		r     gin.H
		err   error
		total int64
	)
	selectStr := "departments.id,departments.name,departments.short_name,departments.sort,departments.headcount,(select count (1) from personnels where personnels.organ_id = departments.id and personnels.status = 1) use"
	result := db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Order("sort desc,level_code asc").Find(&h)
	if err = db.Model(&models.Personnel{}).Where("status = ?", 1).Count(&total).Error; err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	err = result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "data": &h, "total": total}
	c.JSON(200, r)
}

//"select personnels.organ_id, count(case when levels.name = '正科级' and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) zk from posts \nleft join levels on posts.level_id = levels.id \nleft join personnels on posts.personnel_id = personnels.id\ngroup by personnels.organ_id;"

// DepartmentPosition 职数占用
func DepartmentPosition(c *gin.Context) {
	var (
		d   []models.Department
		r   gin.H
		err error
		dp  []struct {
			OrganId string `json:"organId"`
			Zk      int64  `json:"zk"` //正科
			Fk      int64  `json:"fk"` //副科
			Zc      int64  `json:"zc"` //正处
			Fc      int64  `json:"fc"` //副处
			Ft      int64  `json:"ft"` //副厅
		}
	)

	selectStr := "id,name,short_name,sort,position"
	result := db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Order("sort desc,level_code asc").Find(&d)
	err = result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	selectStr = "personnels.organ_id, count(case when levels.name = '正科级' and positions.is_leader = 2 and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) zk, " +
		"count(case when levels.name = '副科级' and positions.is_leader = 2 and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) fk, " +
		"count(case when levels.name = '正处级' and positions.is_leader = 2 and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) zc, " +
		"count(case when levels.name = '副处级' and positions.is_leader = 2 and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) fc, " +
		"count(case when levels.name = '副厅级' and positions.is_leader = 2 and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) ft"
	JoinStr := "left join levels on posts.level_id = levels.id " +
		"left join positions on posts.position_id = positions.id " +
		"left join personnels on posts.personnel_id = personnels.id"
	db.Table("posts").Select(selectStr).Joins(JoinStr).Group("personnels.organ_id").Find(&dp)
	r = gin.H{"code": 20000, "data": &d, "position": &dp}
	c.JSON(200, r)
}

func DepartmentUpdate(c *gin.Context) {
	var h struct {
		ID        string `json:"id"`
		Headcount int    `json:"headcount"`
		Position  string `json:"position"`
	}
	var r gin.H
	field := c.Query("field")
	if c.ShouldBindJSON(&h) != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		c.JSON(200, r)
		return
	}
	if field == "position" {
		db.Table("departments").Where("id = ?", h.ID).Update("position", h.Position)
	} else {
		db.Table("departments").Where("id = ?", h.ID).Update("headcount", h.Headcount)
	}
	r = gin.H{"message": "更新成功！", "code": 20000}
	c.JSON(200, r)
}

// DepartmentHighLevel 高级警长职数占用
func DepartmentHighLevel(c *gin.Context) {
	var (
		d   []models.Department
		r   gin.H
		err error
		dp  []struct {
			OrganId string `json:"organId"`
			G1      int64  `json:"g1"` //一级高级警长
			G2      int64  `json:"g2"` //二级高级警长
			G3      int64  `json:"g3"` //三级高级警长
			G4      int64  `json:"g4"` //四级高级警长
		}
	)

	selectStr := "id,name,short_name,sort,position"
	result := db.Table("departments").Select(selectStr).Where("dept_type = ?", 1).Order("sort desc,level_code asc").Find(&d)
	err = result.Error
	if err != nil {
		//r = Errors.ServerError
		r = GetError(CodeServer)
		c.JSON(200, r)
		return
	}
	selectStr = "personnels.organ_id, count(case when positions.name = '一级高级警长' and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) g1, " +
		"count(case when positions.name = '二级高级警长' and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) g2, " +
		"count(case when positions.name = '三级高级警长' and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) g3, " +
		"count(case when positions.name = '四级高级警长' and posts.end_day = '0001-01-01 00:00:00.000000 +00:00' then 1 else null end) g4"
	JoinStr := "left join positions on posts.position_id = positions.id " +
		"left join personnels on posts.personnel_id = personnels.id"
	db.Table("posts").Select(selectStr).Joins(JoinStr).Group("personnels.organ_id").Find(&dp)
	r = gin.H{"code": 20000, "data": &d, "position": &dp}
	c.JSON(200, r)
}
