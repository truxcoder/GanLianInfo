package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
	log "github.com/truxcoder/truxlog"
)

type OrderStruct struct {
	Id    int64 `json:"id,string"`
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
	db.Table("casbin_rule").Select(selectStr).Where("v1 in ? and v2 = ?", moduleList, "MENU").Find(&mr)
	r = gin.H{"code": 20000, "data": &mr}
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
