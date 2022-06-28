package controllers

import (
	"GanLianInfo/models"
	log "github.com/truxcoder/truxlog"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Awards struct {
	models.Award
	PersonnelName  string `json:"personnelName"`
	PoliceCode     string `json:"policeCode"`
	OrganName      string `json:"organName"`
	OrganShortName string `json:"organShortName"`
}

func AwardList(c *gin.Context) {
	var mos []Awards
	var mo struct {
		models.Award
		Intercept string `json:"intercept" gorm:"-" sql:""`
	}
	selectStr := "awards.*,per.name as personnel_name, per.police_code as police_code," +
		"departments.name as organ_name, departments.short_name as organ_short_name "
	joinStr := "left join personnels as per on awards.personnel_id = per.id " +
		"left join departments on departments.id = per.organ_id "
	getList(c, "awards", &mo, &mos, &selectStr, &joinStr)
}

func AwardDetail(c *gin.Context) {
	var mos []models.Award
	var selectStr string
	var joinStr string
	getDetail(c, "awards", &mos, &selectStr, &joinStr)
}

// AwardPreBatch 批量录入奖励信息之前验证将录入的信息在数据库里是否存在。返还已存在的条目到前端
func AwardPreBatch(c *gin.Context) {
	type temp struct {
		ID          int64 `json:"id,string"`
		PersonnelId int64 `json:"personnelId,string"`
	}
	var (
		r      gin.H
		err    error
		mos    []models.Award
		ids    []int64
		result []temp
	)

	var mo struct {
		Personnels []string  `json:"personnels" gorm:"-"`
		GetTime    time.Time `json:"getTime"`
		Grade      int8      `json:"grade"`
		DocNumber  string    `json:"docNumber"`
	}
	if err = c.ShouldBindJSON(&mo); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	if len(mo.Personnels) > 0 {
		for _, v := range mo.Personnels {
			_v, _ := strconv.Atoi(v)
			ids = append(ids, int64(_v))
		}
	}
	// 查找数据库里是否有指定人员，指定时间, 指定文号的奖励信息
	db.Table("awards").Where(&mo).Where("personnel_id in ?", ids).Find(&mos)
	if len(mos) > 0 {
		for _, v := range mos {
			result = append(result, temp{ID: v.ID, PersonnelId: v.PersonnelId})
		}
	}
	r = gin.H{"code": 20000, "data": result}
	c.JSON(200, r)
}

// AwardBatch 考核信息批量录入
func AwardBatch(c *gin.Context) {
	var (
		r       gin.H
		err     error
		added   []models.Award
		updated []int64
	)

	var mo struct {
		Added     []string  `json:"added" gorm:"-"`   //增加的人员ID列表
		Updated   []string  `json:"updated" gorm:"-"` //修改的考核信息ID列表，注意并不是人员ID列表
		Category  int8      `json:"category"`
		GetTime   time.Time `json:"getTime"`
		Grade     int8      `json:"grade"`
		Content   string    `json:"content"`
		DocNumber string    `json:"docNumber"`
		Organ     string    `json:"organ"`
	}
	if err = c.ShouldBindJSON(&mo); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	if len(mo.Added) > 0 {
		for _, v := range mo.Added {
			_v, _ := strconv.Atoi(v)
			added = append(added, models.Award{PersonnelId: int64(_v),
				Category: mo.Category, GetTime: mo.GetTime, Grade: mo.Grade, Content: mo.Content, DocNumber: mo.DocNumber, Organ: mo.Organ})
		}
	}
	if len(mo.Updated) > 0 {
		for _, v := range mo.Updated {
			_v, _ := strconv.Atoi(v)
			updated = append(updated, int64(_v))
		}
	}
	// 启用事务，确保添加和修改都顺利执行
	err = db.Transaction(func(tx *gorm.DB) error {
		if len(mo.Added) > 0 {
			if _err := tx.Create(added).Error; _err != nil {
				// 返回任何错误都会回滚事务
				return _err
			}
		}
		if len(mo.Updated) > 0 {
			if _err := tx.Table("awards").Where("id in ?", updated).Updates(mo).Error; _err != nil {
				return _err
			}
		}
		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		r = GetError(CodeDatabase)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "message": "操作成功"}
	c.JSON(200, r)
}
