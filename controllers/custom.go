package controllers

import (
	"GanLianInfo/models"

	"github.com/gin-gonic/gin"
)

func CustomList(c *gin.Context) {
	var mos []models.Custom
	var mo struct {
		AccountId string `json:"accountId"`
		Category  int8   `json:"category"`
	}
	selectStr := "*"
	joinStr := ""
	getList(c, "customs", &mo, &mos, &selectStr, &joinStr)
}
