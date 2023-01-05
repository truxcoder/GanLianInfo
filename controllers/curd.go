package controllers

import (
	"GanLianInfo/models"
	"fmt"
	"go/ast"
	"gorm.io/gorm"
	"reflect"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

var maps = make(map[string]reflect.Type)

func init() {
	maps["personnel"] = reflect.TypeOf(models.Personnel{})
	maps["level"] = reflect.TypeOf(models.Level{})
	maps["appraisal"] = reflect.TypeOf(models.Appraisal{})
	maps["post"] = reflect.TypeOf(models.Post{})
	maps["position"] = reflect.TypeOf(models.Position{})
	maps["award"] = reflect.TypeOf(models.Award{})
	maps["punish"] = reflect.TypeOf(models.Punish{})
	maps["module"] = reflect.TypeOf(models.Module{})
	maps["discipline"] = reflect.TypeOf(models.Discipline{})
	maps["training"] = reflect.TypeOf(models.Training{})
	maps["person_train"] = reflect.TypeOf(models.PersonTrain{})
	maps["role_dict"] = reflect.TypeOf(models.RoleDict{})
	maps["permission_dict"] = reflect.TypeOf(models.PermissionDict{})
	maps["dis_dict"] = reflect.TypeOf(models.DisDict{})
	maps["edu_dict"] = reflect.TypeOf(models.EduDict{})
	maps["report"] = reflect.TypeOf(models.Report{})
	maps["entry_exit"] = reflect.TypeOf(models.EntryExit{})
	maps["affair"] = reflect.TypeOf(models.Affair{})
	maps["family"] = reflect.TypeOf(models.Family{})
	maps["talent"] = reflect.TypeOf(models.Talent{})
	maps["custom"] = reflect.TypeOf(models.Custom{})
	maps["review"] = reflect.TypeOf(models.Review{})
	maps["feedback"] = reflect.TypeOf(models.Feedback{})
	maps["appointment"] = reflect.TypeOf(models.Appointment{})
}

func Add(c *gin.Context) {
	var (
		ok bool
		r  gin.H
	)

	model, err := getInstance(c)
	if err != nil {
		log.Error(err)
		return
	}
	if err = c.ShouldBindJSON(model); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	if ok, err = validate(model); err != nil {
		r = GetError(CodeValidate)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	if !ok {
		r = GetError(CodeExist)
		c.JSON(200, r)
		return
	}

	db.Create(model)
	r = gin.H{"message": "添加成功！", "code": 20000}

	c.JSON(200, r)
}

// 添加前数据验证
func validate(model interface{}) (bool, error) {
	var count int64
	if reflect.TypeOf(model) == reflect.TypeOf(&models.Appraisal{}) {
		if mo, ok := model.(*models.Appraisal); ok {
			if err := db.Model(model).Where("personnel_id = ? AND years = ? AND season = ?", mo.PersonnelId, mo.Years, mo.Season).Count(&count).Error; err != nil {
				return false, err
			}
			if count > 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

func Update(c *gin.Context) {
	model, err := getInstance(c)
	if err != nil {
		log.Error(err)
		return
	}
	var r gin.H
	if err = c.ShouldBindJSON(model); err != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		log.Error(err)
	} else {
		db.Model(model).Updates(model)
		updateZeroFields(model)
		r = gin.H{"message": "更新成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func Delete(c *gin.Context) {
	model, err := getInstance(c)
	if err != nil {
		log.Error(err)
		return
	}
	var id IdStruct
	var r gin.H
	if err = c.ShouldBindJSON(&id); err != nil {
		//r = Errors.ServerError
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	// 启用事务，确保所有操作都顺利执行
	err = db.Transaction(func(tx *gorm.DB) error {
		// 如果提交删除的是任职信息
		if reflect.TypeOf(model) == reflect.TypeOf(&models.Post{}) {
			if _err := parsePostDataBeforeDelete(tx, &id); _err != nil {
				return _err
			}
		}

		if _err := tx.Delete(model, &id.Id).Error; _err != nil {
			// 返回任何错误都会回滚事务
			return _err
		}
		// 返回 nil 提交事务
		return nil
	})

	if err != nil {
		r = gin.H{"message": err.Error(), "code": 50500}
		log.Error(err)
		c.JSON(200, r)
		return
	}
	message := fmt.Sprintf("成功删除%d条数据", len(id.Id))
	r = gin.H{"message": message, "code": 20000}
	c.JSON(200, r)

	//result := db.Delete(model, &id.Id)
	//err = result.Error
	//if err != nil {
	//	r = GetError(CodeServer)
	//} else {
	//	message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
	//	r = gin.H{"message": message, "code": 20000}
	//}
	//c.JSON(200, r)
	//return
}

// 在删除任职信息前用事务同步处理当前职务和当前职级。逻辑是：
// 循环要删除的人员id，判断要删除的信息的结束日期是否为零值。如果不是，则不做处理，如果是，则判断是职务还是职级
// 如果是职务，则判断该人员是否还有未结束任期的职务，有则把current_level改为未结束任期的最高职务，无则把current_level置为null
// 如果是职级，则直接将current_rank置为null
func parsePostDataBeforeDelete(tx *gorm.DB, i *IdStruct) error {
	var (
		post      models.Post
		existPost []PostWithLevel
		position  models.Position
	)
	zeroDate := "0001-01-01 00:00:00.000000 +00:00"

	// 循环要删除的人员id，判断要删除的信息的结束日期是否为零值。如果不是，则不做处理，如果是，则同步处理人员的current_level和current_rank
	for _, id := range i.Id {
		if err := tx.Limit(1).Find(&post, id).Error; err != nil {
			return err
		}
		if post.EndDay.IsZero() {
			if err := tx.Limit(1).Find(&position, post.PositionId).Error; err != nil {
				return err
			}
			// 如果提交的是领导职务
			if position.IsLeader == 2 {
				//判断该人员是否还有未结束任期的职务，有则把current_level改为未结束任期的最高职务，无则把current_level置为null
				if err := tx.Model(&models.Post{}).Select("posts.*, levels.name level_name, levels.orders level_order").Joins("left join levels on levels.id = posts.level_id").Where("personnel_id = ? and end_day = ? and position_id in (select id from positions where is_leader = 2)", post.PersonnelId, zeroDate).Find(&existPost).Error; err != nil {
					return err
				}
				if len(existPost) == 1 {
					if err := tx.Model(&models.Personnel{}).Where("id = ?", post.PersonnelId).Update("current_level", nil).Error; err != nil {
						return err
					}
				} else {
					_order := 100
					var _levelId int64
					// 循环判断出级别最高任职信息，取出其id
					for _, v := range existPost {
						if v.LevelOrder < _order && v.ID != id {
							_levelId = v.LevelId
						}
					}
					if _levelId != 0 {
						if err := tx.Model(&models.Personnel{}).Where("id = ?", post.PersonnelId).Update("current_level", _levelId).Error; err != nil {
							return err
						}
					}
				}
				// 如果提交的是非领导职务，则直接将current_rank置为null
			} else if position.IsLeader == 1 {
				if err := tx.Model(&models.Personnel{}).Where("id = ?", post.PersonnelId).Update("current_rank", nil).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// PreEdit 提交审核
func PreEdit(c *gin.Context) {
	var (
		r   gin.H
		mo  models.Review
		err error
	)

	if err = c.ShouldBindJSON(&mo); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	db.Create(&mo)
	r = gin.H{"message": "提交成功！请等待审核", "code": 20000}
	c.JSON(200, r)
}

// PreBatchEdit 提交批量审核
func PreBatchEdit(c *gin.Context) {
	var (
		r   gin.H
		mo  []models.Review
		err error
	)

	if err = c.ShouldBindJSON(&mo); err != nil {
		r = GetError(CodeBind)
		log.Error(err)
		c.JSON(200, r)
		return
	}

	db.Create(&mo)
	r = gin.H{"message": "提交成功！请等待审核", "code": 20000}
	c.JSON(200, r)
}

// getInstance 获取模型实例，返回实例指针
func getInstance(c *gin.Context) (interface{}, error) {
	resource := c.Query("resource")
	t, ok := maps[resource]
	if !ok {
		message := fmt.Sprintf("未找到对应模型: %s", resource)
		err := errors.New(message)
		return nil, err
	}
	return reflect.New(t).Interface(), nil
}

// updateZeroFields 更新零值字段。
// Gorm的Updates或update如果传入的是struct，则忽略零值字段。为了保证一些字段可以归零，在model定义时加上
// update标签，从而构建一个map，用map来更新
func updateZeroFields(model interface{}) {
	var result = make(map[string]interface{})
	var total = 0
	T := reflect.Indirect(reflect.ValueOf(model)).Type()
	V := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < T.NumField(); i++ {
		p := T.Field(i)
		v := V.Field(i)

		if !p.Anonymous && ast.IsExported(p.Name) {
			if !v.IsZero() {
				continue
			}
			if tag, ok := p.Tag.Lookup("update"); ok {
				result[tag] = v.Interface()
				total++
			}
		}
	}
	if total > 0 {
		db.Model(model).Updates(result)
		//log.Successf("更新了%d个零值字段:%+v\n", total, result)
	}
}

// 把零值更新字段转为map
func updateZeroFieldsToMap(model interface{}) (map[string]interface{}, int) {
	var result = make(map[string]interface{})
	var total = 0
	T := reflect.Indirect(reflect.ValueOf(model)).Type()
	V := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < T.NumField(); i++ {
		p := T.Field(i)
		v := V.Field(i)

		if !p.Anonymous && ast.IsExported(p.Name) {
			if !v.IsZero() {
				continue
			}
			if tag, ok := p.Tag.Lookup("update"); ok {
				result[tag] = v.Interface()
				total++
			}
		}
	}
	return result, total
}
