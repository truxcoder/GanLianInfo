package controllers

import (
	"GanLianInfo/models"
	"strconv"

	"gorm.io/gorm"

	log "github.com/truxcoder/truxlog"

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
	var mo struct {
		PersonnelId int64    `json:"personnelId,string"`
		Years       []string `json:"years" gorm:"-" query:"appraisals.years"`
		Season      []int8   `json:"season" gorm:"-" query:"appraisals.season"`
		Conclusion  string   `json:"conclusion"`
	}
	selectStr := "appraisals.*,per.name as personnel_name, per.police_code as police_code," +
		"departments.name as organ_name, departments.short_name as organ_short_name "
	joinStr := "left join personnels as per on appraisals.personnel_id = per.id " +
		"left join departments on departments.id = per.organ_id "
	getList(c, "appraisals", &mo, &mos, &selectStr, &joinStr)

}

func AppraisalDetail(c *gin.Context) {
	//var mos []AppOrgan
	var mos []models.Appraisal
	var selectStr, joinStr string
	//selectStr := "appraisals.*,departments.name as organ_name, departments.short_name as organ_short_name "
	//joinStr := "left join personnels as per on appraisals.personnel_id = per.id " +
	//	"left join departments on departments.id = per.organ_id "
	//getDetail(c, "appraisals", &mos, &selectStr, &joinStr)
	getDetail(c, "appraisals", &mos, &selectStr, &joinStr)
}

// AppraisalPreBatch 批量录入考核信息之前验证将录入的信息在数据库里是否存在。返还已存在的条目到前端
func AppraisalPreBatch(c *gin.Context) {
	type temp struct {
		ID          int64 `json:"id,string"`
		PersonnelId int64 `json:"personnelId,string"`
	}
	var (
		r      gin.H
		err    error
		mos    []models.Appraisal
		ids    []int64
		result []temp
	)

	var mo struct {
		Personnels []string `json:"personnels" gorm:"-"`
		Years      string   `json:"years"`
		Season     int8     `json:"season"`
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
	// 查找数据库里是否有指定人员，指定时间的考核信息
	db.Table("appraisals").Where(&mo).Where("personnel_id in ?", ids).Find(&mos)
	if len(mos) > 0 {
		for _, v := range mos {
			result = append(result, temp{ID: v.ID, PersonnelId: v.PersonnelId})
		}
	}
	r = gin.H{"code": 20000, "data": result}
	c.JSON(200, r)
}

// AppraisalBatch 考核信息批量录入
func AppraisalBatch(c *gin.Context) {
	var (
		r       gin.H
		err     error
		added   []models.Appraisal
		updated []int64
	)

	var mo struct {
		Added      []string `json:"added" gorm:"-"`   //增加的人员ID列表
		Updated    []string `json:"updated" gorm:"-"` //修改的考核信息ID列表，注意并不是人员ID列表
		Organ      string   `json:"organ"`
		Years      string   `json:"years"`
		Season     int8     `json:"season"`
		Conclusion string   `json:"conclusion"`
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
			added = append(added, models.Appraisal{PersonnelId: int64(_v),
				Organ: mo.Organ, Years: mo.Years, Season: mo.Season, Conclusion: mo.Conclusion})
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
			if _err := tx.Table("appraisals").Where("id in ?", updated).Update("conclusion", mo.Conclusion).Error; _err != nil {
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
