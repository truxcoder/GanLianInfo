package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

func FamilyDetail(c *gin.Context) {
	var mos []models.Family
	var selectStr string
	var joinStr string
	getDetail(c, "families", &mos, &selectStr, &joinStr)
}
