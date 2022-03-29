package controllers

import (
	"GanLianInfo/models"
	"fmt"
	"go/ast"
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
}

func Add(c *gin.Context) {
	model, err := getInstance(c)
	if err != nil {
		log.Error(err)
		return
	}
	var r gin.H
	if err = c.ShouldBindJSON(model); err != nil {
		r = Errors.ServerError
		log.Error(err)
	} else {
		db.Create(model)
		r = gin.H{"message": "添加成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func Update(c *gin.Context) {
	model, err := getInstance(c)
	if err != nil {
		log.Error(err)
		return
	}
	var r gin.H
	if err = c.ShouldBindJSON(model); err != nil {
		r = Errors.ServerError
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
		r = Errors.ServerError
		log.Error(err)
		c.JSON(200, r)
		return
	}

	result := db.Delete(model, &id.Id)
	err = result.Error
	if err != nil {
		log.Error(err)
		r = Errors.ServerError
	} else {
		message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
		r = gin.H{"message": message, "code": 20000}
	}
	c.JSON(200, r)
	return
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
