package controllers

import (
	"GanLianInfo/models"
	"GanLianInfo/utils"
	"fmt"
	"reflect"

	"gorm.io/gorm"

	jsoniter "github.com/json-iterator/go"
	log "github.com/truxcoder/truxlog"

	"github.com/gin-gonic/gin"
)

func ReviewList(c *gin.Context) {
	var mos []struct {
		models.Review
		PersonnelName  string `json:"personnelName"`
		ReviewerName   string `json:"reviewerName"`
		OrganShortName string `json:"organShortName"`
	}
	var mo models.Review
	selectStr := "reviews.*,personnels.name as personnel_name,per.name as reviewer_name," +
		"departments.short_name as organ_short_name"
	joinStr := "left join personnels on reviews.personnel_id = personnels.id" +
		" left join personnels as per on reviews.reviewer = per.id" +
		" left join departments on departments.id = reviews.organ_id"
	getList(c, "reviews", &mo, &mos, &selectStr, &joinStr)
}

func ReviewPass(c *gin.Context) {
	var (
		mo          models.Review
		err         error
		r           gin.H
		p           models.Personnel
		edu         PerEduStruct
		resume, res models.Resume
		f           models.Family
		_map        = make(map[string]interface{})
		reviewMap   = make(map[string]interface{})
		dataMap     = make(map[string]interface{})
	)
	if err = c.ShouldBindJSON(&mo); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	reviewMap["status"] = mo.Status
	reviewMap["reviewer"] = mo.Reviewer
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	// 分类为人员基本情况
	if mo.Category == 1 {
		if err = json.Unmarshal([]byte(mo.Content), &_map); err != nil {
			r = GetError(CodeParse)
			log.Error(err)
			c.JSON(200, r)
			return
		}
		if err = json.Unmarshal([]byte(mo.Content), &p); err != nil {
			r = GetError(CodeParse)
			log.Error(err)
			c.JSON(200, r)
			return
		}
		if dataMap, err = bindDataToMap(&p, _map); err != nil {
			r = GetError(CodeParse)
			log.Error(err)
			c.JSON(200, r)
			return
		}

		log.Successf("dataMap-----------:\n %+v\n", dataMap)
		log.Successf("mo-----------:\n %+v\n", mo)
		// 启用事务，确保所有操作都顺利执行
		err = db.Transaction(func(tx *gorm.DB) error {
			if _err := tx.Model(&models.Review{}).Where("id = ?", mo.ID).Updates(reviewMap).Error; _err != nil {
				// 返回任何错误都会回滚事务
				return _err
			}
			if _err := tx.Model(&models.Personnel{}).Where("id = ?", mo.PersonnelId).Updates(dataMap).Error; _err != nil {
				return _err
			}
			// 返回 nil 提交事务
			return nil
		})
	}
	// 分类为教育情况
	if mo.Category == 2 {
		if err = json.Unmarshal([]byte(mo.Content), &edu); err != nil {
			r = GetError(CodeParse)
			log.Error(err)
			c.JSON(200, r)
			return
		}
		if _map, err = utils.StructToMap(&edu); err != nil {
			r = GetError(CodeBind)
			log.Error(err)
			c.JSON(200, r)
			return
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			if _err := tx.Model(&models.Review{}).Where("id = ?", mo.ID).Updates(reviewMap).Error; _err != nil {
				return _err
			}
			if _err := tx.Model(&models.Personnel{}).Where("id = ?", mo.PersonnelId).Updates(_map).Error; _err != nil {
				return _err
			}
			return nil
		})
	}

	//分类为个人简历
	if mo.Category == 3 {
		resume.PersonnelId = mo.PersonnelId
		resume.Content = string(mo.Content)
		result := db.Table("resumes").Where("personnel_id = ?", mo.PersonnelId).Limit(1).Find(&res)
		if result.RowsAffected > 0 {
			err = db.Transaction(func(tx *gorm.DB) error {
				if _err := tx.Model(&models.Review{}).Where("id = ?", mo.ID).Updates(reviewMap).Error; _err != nil {
					return _err
				}
				if _err := tx.Model(&models.Resume{}).Where("id = ?", res.ID).Updates(&resume).Error; _err != nil {
					return _err
				}
				return nil
			})
		} else {
			err = db.Transaction(func(tx *gorm.DB) error {
				if _err := tx.Model(&models.Review{}).Where("id = ?", mo.ID).Updates(reviewMap).Error; _err != nil {
					return _err
				}
				if _err := tx.Model(&models.Resume{}).Create(&resume).Error; _err != nil {
					return _err
				}
				return nil
			})
		}
	}

	// 分类为家庭成员
	if mo.Category == 4 {
		if err = json.Unmarshal([]byte(mo.Content), &f); err != nil {
			r = GetError(CodeParse)
			log.Error(err)
			c.JSON(200, r)
			return
		}
		err = db.Transaction(func(tx *gorm.DB) error {
			if _err := tx.Model(&models.Review{}).Where("id = ?", mo.ID).Updates(reviewMap).Error; _err != nil {
				return _err
			}
			if f.ID == 0 {
				if _err := tx.Create(&f).Error; _err != nil {
					return _err
				}
			} else if _err := tx.Model(&f).Where("id = ?", &f.ID).Updates(&f).Error; _err != nil {
				return _err
			}
			return nil
		})
	}

	if err != nil {
		r = GetError(CodeDatabase)
		log.Error(err)
		c.JSON(200, r)
		return
	}
	r = gin.H{"code": 20000, "message": "操作成功"}
	c.JSON(200, r)
}

func bindDataToMap(in interface{}, m map[string]interface{}) (map[string]interface{}, error) {
	// 当前函数只接收struct类型
	inV := reflect.Indirect(reflect.ValueOf(in))
	inT := reflect.Indirect(reflect.ValueOf(in)).Type()
	if inT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToMap函数的参数只能为struct指针; got %+v", inT)
	}

	out := make(map[string]interface{})
	for i := 0; i < inT.NumField(); i++ {
		p := inT.Field(i)
		if !p.Anonymous {
			if _, ok := m[p.Tag.Get("json")]; ok {
				out[p.Name] = inV.Field(i).Interface()
			}
		} else {
			field := inV.Field(i)
			for j := 0; j < p.Type.NumField(); j++ {
				pp := p.Type.Field(j)
				if _, ok := m[pp.Tag.Get("json")]; ok {
					out[pp.Name] = field.Field(j).Interface()
				}
			}
		}
	}
	return out, nil
}
