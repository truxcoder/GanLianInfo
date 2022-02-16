package controllers

import (
	"GanLianInfo/models"
	"fmt"
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
		log.Successf("model:%+v\n", model)
		db.Debug().Model(model).Updates(model)
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
