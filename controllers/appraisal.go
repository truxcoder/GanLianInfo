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
	var mo AppPerson
	selectStr := "appraisals.*,d.name as organ_name,d.short_name as organ_short_name," +
		"p.name as personnel_name,p.police_code as police_code"
	joinStr := "left join departments as d on appraisals.organ_id = d.id " +
		"left join personnels as p on appraisals.personnel_id = p.id"
	getList(c, "appraisals", &mo, &mos, &selectStr, &joinStr)

}

func AppraisalDetail(c *gin.Context) {
	var mos []AppOrgan
	selectStr := "appraisals.*,d.name as organ_name,d.short_name as organ_short_name"
	joinStr := "left join departments as d on appraisals.organ_id = d.id "
	getDetail(c, "appraisals", &mos, &selectStr, &joinStr)
}

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
		OrganId    string   `json:"organId"`
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
			//mos = append(mos, models.Appraisal{PersonnelId: int64(_v),
			//	OrganId: mo.OrganId, Years: mo.Years, Season: mo.Season, Conclusion: mo.Conclusion})
		}
	}
	db.Table("appraisals").Where(&mo).Where("personnel_id in ?", ids).Find(&mos)
	if len(mos) > 0 {
		for _, v := range mos {
			result = append(result, temp{ID: v.ID, PersonnelId: v.PersonnelId})
		}
	}
	r = gin.H{"code": 20000, "data": result}
	c.JSON(200, r)
}

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
		OrganId    string   `json:"organId"`
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
				OrganId: mo.OrganId, Years: mo.Years, Season: mo.Season, Conclusion: mo.Conclusion})
		}
	}
	if len(mo.Updated) > 0 {
		for _, v := range mo.Updated {
			_v, _ := strconv.Atoi(v)
			updated = append(updated, int64(_v))
		}
	}
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
