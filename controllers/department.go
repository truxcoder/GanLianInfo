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

func GetDepartmentList(c *gin.Context) {
	var d []models.Department
	var r gin.H
	result := db.Order("sort desc").Find(&d)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &d}
	}
	c.JSON(200, r)
}

func GetOrganList(c *gin.Context) {
	var o []Organ
	var r gin.H
	result := db.Model(&models.Department{}).Where("dept_type = ?", 1).Order("sort asc").Find(&o)
	err := result.Error
	if err != nil {
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &o}
	}
	c.JSON(200, r)
}

//func DepartmentUpdate(c *gin.Context) {
//	db := dao.Connect()
//	var o models.Organ
//	var r gin.H
//	if c.ShouldBindJSON(&o) != nil {
//		r = Errors.ServerError
//	} else {
//		db.Model(&o).Updates(&o)
//		r = gin.H{"message": "更新成功！", "code": 20000}
//	}
//	c.JSON(200, r)
//}
