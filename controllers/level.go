package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

func LevelList(c *gin.Context) {
	var levels []models.Level
	var r gin.H
	result := db.Order("orders asc").Find(&levels)
	err := result.Error
	if err != nil {
		r = GetError(CodeServer)
	} else {
		r = gin.H{"code": 20000, "data": &levels}
	}
	c.JSON(200, r)
}
