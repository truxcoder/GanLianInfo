package controllers

import (
	"GanLianInfo/models"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type OrderStruct struct {
	Id    int64 `json:"id"`
	Order int8  `json:"order"`
}

func GetModuleList(c *gin.Context) {
	var m []models.Module
	db.Order("orders asc").Find(&m)
	r := gin.H{"code": 20000, "data": &m}
	c.JSON(200, r)
}

func ModuleRole(c *gin.Context) {
	var moduleList []string
	var err error
	var r gin.H
	var mr []struct {
		Module string `json:"module"`
		Roles  string `json:"roles"`
	}
	if err = c.BindJSON(&moduleList); err != nil {
		log.Error(err)
		r = Errors.ServerError
		c.JSON(200, r)
		return
	}
	selectStr := "v1 as module,v0 as roles"
	db.Debug().Table("casbin_rule").Select(selectStr).Where("v1 in ? and v2 = ?", moduleList, "MENU").Find(&mr)
	r = gin.H{"code": 20000, "data": &mr}
	c.JSON(200, r)
}

func ModuleAdd(c *gin.Context) {
	var m models.Module
	var r gin.H
	err := c.ShouldBindJSON(&m)
	if err != nil {
		r = Errors.ServerError
		log.Error(err)
	} else {
		db.Create(&m)
		r = gin.H{"message": "添加成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func ModuleUpdate(c *gin.Context) {
	var m models.Module
	var r gin.H
	if err := c.ShouldBindJSON(&m); err != nil {
		r = Errors.ServerError
		log.Error(err)
	} else {
		db.Model(&m).Updates(&m)
		r = gin.H{"message": "更新成功！", "code": 20000}
	}
	c.JSON(200, r)
}

func ModuleDelete(c *gin.Context) {
	var id IdStruct
	var r gin.H
	if err := c.ShouldBindJSON(&id); err != nil {
		r = Errors.ServerError
		log.Error(err)
	} else {
		//result := db.Where(&id.Id).Delete(&models.Department{})
		result := db.Delete(&models.Module{}, &id.Id)
		err := result.Error
		if err != nil {
			log.Error(err)
			r = Errors.ServerError
		} else {
			message := fmt.Sprintf("成功删除%d条数据", result.RowsAffected)
			r = gin.H{"message": message, "code": 20000}
		}
	}
	c.JSON(200, r)
}

func ModuleOrder(c *gin.Context) {
	var orders []OrderStruct
	var r gin.H
	if err := c.ShouldBindJSON(&orders); err != nil {
		r = Errors.ServerError
		log.Error(err)
	} else {
		for _, v := range orders {
			db.Model(&models.Module{}).Where("id = ?", v.Id).Update("Orders", v.Order)
		}
		r = gin.H{"message": "更新序号成功！", "code": 20000}
	}
	c.JSON(200, r)
}
