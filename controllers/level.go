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
		r = Errors.ServerError
	} else {
		r = gin.H{"code": 20000, "data": &levels}
	}
	//time.Sleep(4 * time.Second)
	c.JSON(200, r)
}
